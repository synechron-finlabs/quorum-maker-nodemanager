package contracthandler

type RequestHandler interface {
	Encode() string
}

type ResponseHandler interface {
	Decode(r string)
}

type ContractParam struct {
	From    string
	To      string
	Passwd  string
	Parties []string
}

type DeployContractHandler struct {
	Binary string
}

func (d DeployContractHandler) Encode() string {

	return d.Binary
}
