package controllerhelpers

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/utils/ptr"
)

type objectForLabelsPatch struct {
	objectMetaForLabelsPatch `json:"metadata"`
}

type objectMetaForLabelsPatch struct {
	ResourceVersion string             `json:"resourceVersion"`
	Labels          map[string]*string `json:"labels"`
}

func PrepareSetLabelsPatch(obj metav1.Object, labels map[string]*string) ([]byte, error) {
	newLabels := make(map[string]*string, len(obj.GetLabels())+len(labels))
	for k, v := range obj.GetLabels() {
		newLabels[k] = ptr.To(v)
	}

	for k, v := range labels {
		newLabels[k] = v
	}

	patch, err := json.Marshal(objectForLabelsPatch{
		objectMetaForLabelsPatch: objectMetaForLabelsPatch{
			ResourceVersion: obj.GetResourceVersion(),
			Labels:          newLabels,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("can't marshal object for set labels patch: %w", err)
	}

	return patch, nil
}
