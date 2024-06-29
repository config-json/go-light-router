package lightrouter

// HTTP methods
type HTTPMethod string

const (
	GET    HTTPMethod = "GET"
	POST   HTTPMethod = "POST"
	PUT    HTTPMethod = "PUT"
	PATCH  HTTPMethod = "PATCH"
	DELETE HTTPMethod = "DELETE"
)

func ContentTypeHeader(fileExtension string) string {
	switch fileExtension {
	case "html":
		return "text/html"
	case "css":
		return "text/css"
	case "js":
		return "text/javascript"
	case "pdf":
		return "application/pdf"
	case "png":
		return "image/png"
	case "jpg":
		return "image/vnd.sealedmedia.softseal.jpg"
	case "ico":
		return "image/x-icon"
	case "docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	default:
		return "application/octet-stream"
	}
}
