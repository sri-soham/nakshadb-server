package helper

type Result struct {
	is_success bool
	is_fatal   bool
	data       map[string]interface{}
	errors     []string
	err        error
}

func (r *Result) IsSuccess() bool {
	return r.is_success
}

func (r *Result) IsFatalError() bool {
	return r.is_fatal
}

func (r *Result) AddToData(key string, value interface{}) {
	r.data[key] = value
}

func (r *Result) AppendError(error string) {
	r.errors = append(r.errors, error)
}

func (r *Result) GetData() map[string]interface{} {
	return r.data
}

func (r *Result) SetErrors(errors []string) {
	r.errors = errors
}

func (r *Result) GetErrors() []string {
	return r.errors
}

func (r *Result) GetDataByKey(key string) (interface{}, bool) {
	i, ok := r.data[key]

	return i, ok
}

func (r *Result) GetIntData(key string) int {
	v, _ := r.data[key]
	return v.(int)
}

func (r *Result) GetInt64Data(key string) int64 {
	v, _ := r.data[key]
	return v.(int64)
}

func (r *Result) GetStringData(key string) string {
	v, _ := r.data[key]
	return v.(string)
}

func (r *Result) GetBoolData(key string) bool {
	v, _ := r.data[key]
	return v.(bool)
}

func (r *Result) GetFloat32Data(key string) float32 {
	v, _ := r.data[key]
	return v.(float32)
}

func (r *Result) GetFloat64Data(key string) float64 {
	v, _ := r.data[key]
	return v.(float64)
}

func (r *Result) GetError() error {
	return r.err
}

func (r *Result) ForJsonResponse() map[string]interface{} {
	data := make(map[string]interface{})

	if r.is_success {
		data = r.data
		data["status"] = "success"
	} else {
		data["status"] = "error"
		if r.is_fatal {
			data["errors"] = r.err
		} else {
			data["errors"] = r.errors
		}
	}

	return data
}

func MakeSuccessResult() Result {
	return Result{true, false, make(map[string]interface{}), make([]string, 0), nil}
}

func MakeErrorResult() Result {
	return Result{false, false, make(map[string]interface{}), make([]string, 0), nil}
}

func MakeFatalResult(err error) Result {
	return Result{false, true, make(map[string]interface{}), make([]string, 0), err}
}

func ErrorResultError(err error) Result {
	return Result{false, false, make(map[string]interface{}), []string{err.Error()}, nil}
}

func ErrorResultString(err_str string) Result {
	return Result{false, false, make(map[string]interface{}), []string{err_str}, nil}
}
