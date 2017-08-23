package workers

import (
	"strings"
	"regexp"
	"net/url"
	"strconv"
)

// CloudFrontLog contains all data for one log line
// Documentation http://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/AccessLogs.html
type CloudFrontLog struct {
	Date                    string  `workers:"date"`
	Time                    string  `workers:"time"`
	XEdgeLocation           string  `workers:"x_edge_location"`
	ScBytes                 int64   `workers:"sc_bytes"`
	CIp                     string  `workers:"c_ip"`
	CsMethod                string  `workers:"cs_method"`
	CsHost                  string  `workers:"cs_host"`
	CsUriStem               string  `workers:"cs_uri_stem"`
	ScStatus                int16   `workers:"sc_status"`
	CsReferer               string  `workers:"cs_referer"`
	CsUserAgent             string  `workers:"cs_user_agent"`
	CsUriQuery              string  `workers:"cs_uri_query"`
	CsCookie                string  `workers:"cs_cookie"`
	XEdgeResultType         string  `workers:"x_edge_result_type"`
	XEdgeRequestId          string  `workers:"x_edge_request_id"`
	XHostHeader             string  `workers:"x_host_header"`
	CsProtocol              string  `workers:"cs_protocol"`
	CsBytes                 int32   `workers:"cs_bytes"`
	TimeTaken               float64 `workers:"time_taken"`
	XForwardedFor           string  `workers:"x_forwarded_for"`
	SslProtocol             string  `workers:"ssl_protocol"`
	SslCipher               string  `workers:"ssl_cipher"`
	XEdgeResponseResultType string  `workers:"x_edge_response_result_type"`
	CsProtocolVersion       string  `workers:"cs_protocol_version"`
}

// CloudFrontLogs is a list of log lines
type CloudFrontLogs struct {
	Logs []CloudFrontLog
}

// SetContent is a string of several log lines
func (c *CloudFrontLogs) SetContent(content string) *CloudFrontLogs {
	// Empty previous content
	for i := 0; i < len(c.Logs); i++ {
		c.Logs[i] = CloudFrontLog{}
	}
	c.Logs = c.Logs[:0]

	// Separate log lines
	lineSlice := strings.Split(content, "\n")

	// Loop trough log lines
	for _, rawLine := range lineSlice {
		line := strings.TrimSpace(rawLine)
		if c.isLogLine(&line) {
			logSplit := strings.Split(line, "\t")
			if len(logSplit) >= 24 {
				c.addLogLine(&logSplit)
			} else {
				//c.logger.Log("To few fields in " + line)
			}
		}
	}

	return c
}

// addLogLine append new log line to the list
func (c *CloudFrontLogs) addLogLine(logLineFields *[]string) {
	c.Logs = append(c.Logs, CloudFrontLog{
		Date:                    c.getDate(logLineFields),
		Time:                    c.getTime(logLineFields),
		XEdgeLocation:           c.getEdgeLocation(logLineFields),
		ScBytes:                 c.getScBytes(logLineFields),
		CIp:                     c.getIpAddress(logLineFields),
		CsMethod:                c.getCsMethod(logLineFields),
		CsHost:                  c.getCsHost(logLineFields),
		CsUriStem:               c.getCsUri(logLineFields),
		ScStatus:                c.getScStatus(logLineFields),
		CsReferer:               c.getReferrer(logLineFields),
		CsUserAgent:             c.getUserAgent(logLineFields),
		CsUriQuery:              c.getCsUriQuery(logLineFields),
		CsCookie:                c.getCookie(logLineFields),
		XEdgeResultType:         c.getEdgeResultType(logLineFields),
		XEdgeRequestId:          c.getEdgeRequestId(logLineFields),
		XHostHeader:             c.getHostHeader(logLineFields),
		CsProtocol:              c.getProtocol(logLineFields),
		CsBytes:                 c.getCsBytes(logLineFields),
		TimeTaken:               c.getProcessTime(logLineFields),
		XForwardedFor:           c.getXForwardedFor(logLineFields),
		SslProtocol:             c.getSSlProtocol(logLineFields),
		SslCipher:               c.getSslCipher(logLineFields),
		XEdgeResponseResultType: c.getResponseType(logLineFields),
		CsProtocolVersion:       c.getProtocolVersion(logLineFields),
	})
}

// isLogLine return true if the line is a log line
func (c *CloudFrontLogs) isLogLine(logLine *string) bool {
	isComment, _ := regexp.MatchString("^#", *logLine)
	if !isComment && *logLine != "" {
		return true
	} else {
		return false
	}
}

func (c *CloudFrontLogs) getDate(log *[]string) string {
	return (*log)[0]
}

func (c *CloudFrontLogs) getTime(log *[]string) string {
	return (*log)[1]
}

func (c *CloudFrontLogs) getEdgeLocation(log *[]string) string {
	return (*log)[2]
}

func (c *CloudFrontLogs) getScBytes(log *[]string) int64 {
	integer, _ := strconv.ParseInt((*log)[3], 10, 64)
	return integer
}

func (c *CloudFrontLogs) getIpAddress(log *[]string) string {

	ip := (*log)[4]

	// Is it a real IP address
	realIP, _ := regexp.MatchString(`[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}`, ip)
	if realIP {
		return ip
	} else {
		return ""
	}
}

func (c *CloudFrontLogs) getCsMethod(log *[]string) string {
	return (*log)[5]
}

func (c *CloudFrontLogs) getCsHost(log *[]string) string {
	return (*log)[6]
}

func (c *CloudFrontLogs) getCsUri(log *[]string) string {
	return (*log)[7]
}

func (c *CloudFrontLogs) getScStatus(log *[]string) int16 {
	integer, _ := strconv.ParseInt((*log)[8], 10, 16)
	return int16(integer)
}

func (c *CloudFrontLogs) getReferrer(log *[]string) string {
	return (*log)[9]
}

func (c *CloudFrontLogs) getUserAgent(log *[]string) string {

	// Url decoding is needed twice
	urlEscape, _ := url.QueryUnescape((*log)[10])
	urlEscape, _ = url.QueryUnescape(urlEscape)

	return string(urlEscape)
}

func (c *CloudFrontLogs) getCsUriQuery(log *[]string) string {
	return (*log)[11]
}

func (c *CloudFrontLogs) getCookie(log *[]string) string {
	return (*log)[12]
}

func (c *CloudFrontLogs) getEdgeResultType(log *[]string) string {
	return (*log)[13]
}

func (c *CloudFrontLogs) getEdgeRequestId(log *[]string) string {
	return (*log)[14]
}

func (c *CloudFrontLogs) getHostHeader(log *[]string) string {
	return (*log)[15]
}

func (c *CloudFrontLogs) getProtocol(log *[]string) string {
	return (*log)[16]
}

func (c *CloudFrontLogs) getCsBytes(log *[]string) int32 {
	integer, _ := strconv.ParseInt((*log)[17], 10, 32)
	return int32(integer)
}

func (c *CloudFrontLogs) getProcessTime(log *[]string) float64 {
	float, _ := strconv.ParseFloat((*log)[18], 64)
	return float
}

func (c *CloudFrontLogs) getXForwardedFor(log *[]string) string {
	return (*log)[19]
}

func (c *CloudFrontLogs) getSSlProtocol(log *[]string) string {
	return (*log)[20]
}

func (c *CloudFrontLogs) getSslCipher(log *[]string) string {
	return (*log)[21]
}

func (c *CloudFrontLogs) getResponseType(log *[]string) string {
	return (*log)[22]
}

func (c *CloudFrontLogs) getProtocolVersion(log *[]string) string {
	return (*log)[23]
}
