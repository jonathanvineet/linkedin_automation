// API Client for LinkedIn Automation Backend

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8090/api';

interface ApiResponse<T = any> {
  success?: boolean;
  error?: string;
  message?: string;
  data?: T;
}

export interface SystemStatus {
  running: boolean;
  logged_in: boolean;
  persona: string;
  stealth: boolean;
  timestamp: string;
}

export interface Stats {
  connections_sent: number;
  messages_sent: number;
  cooldown_seconds: number;
  daily_limit: {
    connections: number;
    messages: number;
  };
}

export interface ActivityLogEntry {
  id: string;
  timestamp: string;
  action: string;
  type: 'info' | 'success' | 'warning' | 'error';
  details: string;
}

export interface SearchRequest {
  keywords?: string;
  location?: string;
  company?: string;
  job_title?: string;
  max_results?: number;
}

export interface ConnectionRequest {
  profile_url: string;
  note?: string;
}

export interface MessageRequest {
  profile_url: string;
  message: string;
}

class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  private async request<T>(
    endpoint: string,
    options?: RequestInit
  ): Promise<T> {
    try {
      const response = await fetch(`${this.baseUrl}${endpoint}`, {
        ...options,
        headers: {
          'Content-Type': 'application/json',
          ...options?.headers,
        },
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || `HTTP ${response.status}`);
      }

      return await response.json();
    } catch (error) {
      console.error(`API request failed: ${endpoint}`, error);
      throw error;
    }
  }

  // Status endpoint
  async getStatus(): Promise<SystemStatus> {
    return this.request<SystemStatus>('/status');
  }

  // Start automation
  async start(): Promise<ApiResponse> {
    return this.request<ApiResponse>('/start', {
      method: 'POST',
    });
  }

  // Stop automation
  async stop(): Promise<ApiResponse> {
    return this.request<ApiResponse>('/stop', {
      method: 'POST',
    });
  }

  // Get statistics
  async getStats(): Promise<Stats> {
    return this.request<Stats>('/stats');
  }

  // Get activity log
  async getActivity(): Promise<ActivityLogEntry[]> {
    return this.request<ActivityLogEntry[]>('/activity');
  }

  // Change persona
  async setPersona(personaType: string): Promise<ApiResponse> {
    return this.request<ApiResponse>('/persona', {
      method: 'POST',
      body: JSON.stringify({ persona_type: personaType }),
    });
  }

  // Search people
  async search(criteria: SearchRequest): Promise<ApiResponse> {
    return this.request<ApiResponse>('/search', {
      method: 'POST',
      body: JSON.stringify(criteria),
    });
  }

  // Send connection request
  async sendConnectionRequest(request: ConnectionRequest): Promise<ApiResponse> {
    return this.request<ApiResponse>('/connect', {
      method: 'POST',
      body: JSON.stringify(request),
    });
  }

  // Send message
  async sendMessage(request: MessageRequest): Promise<ApiResponse> {
    return this.request<ApiResponse>('/message', {
      method: 'POST',
      body: JSON.stringify(request),
    });
  }
}

export const apiClient = new ApiClient(API_BASE_URL);
