package auditmodel

type OperationLog struct {
	UserID        *uint64
	SourceIP      string
	RequestMethod string
	RequestPath   string
	StatusCode    int
	CostMs        int64
	UserAgent     string
	ErrorMessage  string
	RequestBody   string
	ResponseBody  string
}
