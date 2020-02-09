package labels

import (
	"github.com/apirator/apirator/pkg/apis/apirator/v1alpha1"
)

func LabelForAPIMock(mock *v1alpha1.APIMock) map[string]string {
	return map[string]string{"app": "apirator", "api": mock.Name, "managed-by": "apirator"}
}
