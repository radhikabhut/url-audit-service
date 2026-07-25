export interface AuditRequest {
  url: string;
}

export interface AuditResult {
  url: string;
  reachable: boolean;
  statusCode: number;
  responseTimeMs: number;
  contentType: string;
  contentLength: number;
  title: string;
  cached: boolean;
  checkedAt: string;
}

export interface APIErrorDetail {
  code: string;
  message: string;
}

export interface APIError {
  requestId: string;
  error: APIErrorDetail;
}
