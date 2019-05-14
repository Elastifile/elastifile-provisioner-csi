package efaas

import (
	"fmt"
	"testing"
)

func TestAPI(t *testing.T) {
	t.Parallel()

	data := []byte(`
	{
	 "type": "service_account",
	 "project_id": "elastifile-gce-lab-c934",
	 "private_key_id": "5e0d188967e7f23ad77129ff4c9ab59889ccd25d",
	 "private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCMBJyta1PEkd7q\nCLEYNdUBqk4Hlnw7mGXnByjao+4SOZi7mJ1NIAtYjptJ/rcPxjft+hxEba1a1DON\nUU7RuJ3eQk+kLVHdbD2D4noMw6VxJtuWnuyQ2V8v5ojv8kVvVSsbkDAQHVGKTe/8\nCEHxlekGoY0NC+KwWlUKmb7cv/B/2aD1eFsyV7ALE/YJmyFbbvtLrab+U5js04ER\nIWcE+gKlvAF7Xq9Iq6MucyjRvgPagz5RSP146HjbCPdJIz3ilcEL7idVGaZnnx/P\ncZAqYnYZAJTGBhi4fUEpAYR7KVUWIVfc9oXEKJDNwwBHnyyZMBPdYn9prs7xgrEL\ngA+WHPPZAgMBAAECggEACVNhUBee66+/hhzwFqm3NzYtnknCmoGK//k1GmLiv2oA\npzYB/BoPR2WwKByD+tP786i96zzW1/7cNCRfOI6wTRZjkY7HLhVAf6E8+c6qHUA2\nTfDl1rvzoBAdvMWJJGIqzdorqVcakDiirEmsgre2Xo+yAlVxUsehdGRLFw7dqNYv\nrINMqjE2W/SCd8jw2WmplmH+c0MvBKkving9CCNgFnvSMUGinv7y3Zvf2GpplvlC\nFdSFGGXxn1o6HbgrkovKn6EVZ8nP3JadG5evwjotEv1fcEu4vOKMq/jgvfxzscRf\ng9bfdhb3/oc+x43dsH3fR0axaImB7LKKgfu7w7vnJQKBgQDCmgAE7noPd0bt7Xg+\nrl44OgCHv3x0QY4lx0y07Yo1Bg1C72H8BCghr/5rxGUOSCGjoFYTVeLhCVIsYX+8\nxbtplxCJFAgN7lu48EyCgIpP7ppjf1a3Uh762O04BCMw0tXw22ich7d4KN5+r8L7\nOknRStrZYD89QjoUsSEYOK0wnwKBgQC4MePUNoBJEG+yhlMOpDz7mnf/F1U4gFQQ\nxD4stAEA1P/QuSgMb0snJJA3yT3dCL4W2DUxDCWOH/Wx3XnJy216+QR//8fHImCR\nYS4fjmaWlbMOKko1yeCtCLsNfA5uB5Yplrujn2o6v5BE52h3JCjW4qUqzZ6T9cBq\n0rQFacWwhwKBgBKLJDdUFjOFFTA08cFfUkEfXc+RsqVNXeNBs5CGFiZpVjgroXWn\nW7+iCqdwRoTu4K276JfdFkqFXdw2yjpNyUcNixjU3NOfBASCeXfyEbv+K54Rk0zS\nuXsD0s8ErenIHXTfI3/O+u+rTVBbJURVUJVuAZ63Ki+HMQupuVKai/5XAoGBALcp\nHSV8IKsHBhtfSR5JIT8MhoCKIjsyGOYnTrBDOrAqHkveor1iujetOx/OJI80T1oG\nGzavnnSqwTXiR2XrvO1IzDnADletjptiKGxGvSrGp6vRT8QXACzwfpjVIMA3GRI4\nClSVhBvxO7PY7N90fIvaCmX629LD0FgpN8weNu/nAoGAP4rXRr37757Q+c/qeKyU\nsmUCYeHj6w+GIkqJIhsDsj5tE8fLTyU87LF6hvscxYJCX9ZVycvhuzRBiFLkc9yo\nZUKC4SllFDw4Zl63RU7me3PnZHpomiNs0hk3fgqAME1Cx3Pn8NT6iptybSqk2kb7\nHOuPCeblZecVZU0UOPyQrWM=\n-----END PRIVATE KEY-----\n",
	 "client_email": "efaas-csi@elastifile-gce-lab-c934.iam.gserviceaccount.com",
	 "client_id": "102179953128561786237",
	 "auth_uri": "https://accounts.google.com/o/oauth2/auth",
	 "token_uri": "https://oauth2.googleapis.com/token",
	 "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
	 "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/efaas-csi%40elastifile-gce-lab-c934.iam.gserviceaccount.com"
	}
	`)

	res, err := demo1(data)
	if err != nil {
		t.Fatal(fmt.Sprintf("AAAAA %v", err.Error()))
	}

	t.Logf("RES: %v", string(res))
}
