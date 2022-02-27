package forwarder

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"sync"
	"syscall"

	"golang.org/x/sync/errgroup"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

var once sync.Once

func Exec(ctx context.Context, options []*Option, config *restclient.Config) (*Result, error) {
	newOptions, err := parseOptions(options)
	if err != nil {
		return nil, err
	}

	podOptions, err := handleOptions(ctx, newOptions, config)
	if err != nil {
		return nil, err
	}

	stream := genericclioptions.IOStreams{
		In:     os.Stdin,
		Out:    nil,
		ErrOut: os.Stderr,
	}

	carries := make([]*carry, len(podOptions))

	var g errgroup.Group

	for index, option := range podOptions {
		index := index
		stopCh := make(chan struct{}, 1)
		readyCh := make(chan struct{})

		req := &portForwardAPodRequest{
			RestConfig: config,
			Pod:        option.Pod,
			LocalPort:  option.LocalPort,
			PodPort:    option.PodPort,
			Streams:    stream,
			StopCh:     stopCh,
			ReadyCh:    readyCh,
		}
		g.Go(func() error {
			pf, err := portForwardAPod(req)
			if err != nil {
				return err
			}
			carries[index] = &carry{StopCh: stopCh, ReadyCh: readyCh, PF: pf}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Println("wait")
		return nil, err
	}

	ret := &Result{
		Close: func() {
			once.Do(func() {
				for _, c := range carries {
					close(c.StopCh)
				}
			})
		},
		Ready: func() ([][]portforward.ForwardedPort, error) {
			pfs := [][]portforward.ForwardedPort{}
			for _, c := range carries {
				<-c.ReadyCh
				ports, err := c.PF.GetPorts()
				if err != nil {
					return nil, err
				}
				pfs = append(pfs, ports)
			}
			return pfs, nil
		},
	}

	ret.Wait = func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		ret.Close()
	}

	go func() {
		<-ctx.Done()
		ret.Close()
	}()

	return ret, nil
}

func portForwardAPod(req *portForwardAPodRequest) (*portforward.PortForwarder, error) {
	targetURL, err := url.Parse(req.RestConfig.Host)
	if err != nil {
		return nil, err
	}

	targetURL.Path = path.Join(
		"api", "v1",
		"namespaces", req.Pod.Namespace,
		"pods", req.Pod.Name,
		"portforward",
	)

	transport, upgrader, err := spdy.RoundTripperFor(req.RestConfig)
	if err != nil {
		return nil, err
	}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, targetURL)
	fw, err := portforward.New(dialer, []string{fmt.Sprintf("%d:%d", req.LocalPort, req.PodPort)}, req.StopCh, req.ReadyCh, req.Streams.Out, req.Streams.ErrOut)
	if err != nil {
		return nil, err
	}

	go func() {
		if err := fw.ForwardPorts(); err != nil {
			panic(err)
		}
	}()

	return fw, nil
}
