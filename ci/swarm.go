package ci

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/moby/moby/client"
)

func updateSwarmServiceImage(serviceName, image string) error {
	cli, err := client.NewEnvClient()
	if err != nil {
		log().Debugf("Unable to create docker client: %v", err)
		return err
	}
	defer cli.Close()

	service, _, err := cli.ServiceInspectWithRaw(context.Background(), serviceName, types.ServiceInspectOptions{})
	if err != nil {
		log().Debugf("Error inspecting service %s: %v", serviceName, err)
		return err
	}

	_, err = cli.ServiceUpdate(context.Background(), serviceName, service.Meta.Version, swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name: serviceName,
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: swarm.ContainerSpec{
				Image: image,
			},
		},
	}, types.ServiceUpdateOptions{})
	if err != nil {
		log().Debugf("Error updating service %s: %v", serviceName, err)
		return err
	}

	return nil
}
