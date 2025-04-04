type FetchOptions = {
  method?: string;
  headers?: Record<string, string>;
  body?: any;
  requireAuth?: boolean;
};

export class ApiError extends Error {
  status: number;

  constructor(message: string, status: number) {
    super(message);
    this.status = status;
    this.name = 'ApiError';
  }
}

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export const fetchApi = async (
  endpoint: string,
  token: string | null,
  options: FetchOptions = {}
) => {
  const { method = 'GET', headers = {}, body, requireAuth = true } = options;

  const requestHeaders: Record<string, string> = {
    'Content-Type': 'application/json',
    ...headers,
  };

  if (requireAuth) {
    if (!token) {
      throw new ApiError('Authentication required', 401);
    }
    requestHeaders['Authorization'] = `Bearer ${token}`;
  }

  const requestOptions: RequestInit = {
    method,
    headers: requestHeaders,
    body: body ? JSON.stringify(body) : undefined,
  };

  try {
    const response = await fetch(`${API_BASE_URL}${endpoint}`, requestOptions);

    // Parse the JSON response if the content type is application/json
    const contentType = response.headers.get('content-type');
    const data = contentType?.includes('application/json')
      ? await response.json()
      : await response.text();

    if (!response.ok) {
      throw new ApiError(
        typeof data === 'object' && data.message ? data.message : 'API request failed',
        response.status
      );
    }

    return data;
  } catch (error) {
    if (error instanceof ApiError) {
      throw error;
    }
    throw new ApiError((error as Error).message || 'Network error', 500);
  }
};
