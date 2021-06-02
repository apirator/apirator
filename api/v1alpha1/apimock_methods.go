package v1alpha1

func (in *APIMock) MatchLabels() map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       in.GetName(),
		"app.kubernetes.io/managed-by": "apirator",
	}
}
