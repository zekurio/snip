import { AccessTokenModel, ErrorResponse } from "./models";

import { APIError } from "./errors";
import { replaceDoublePath } from "./util";

export type HttpMethod =
  | "GET"
  | "PUT"
  | "POST"
  | "DELETE"
  | "PATCH"
  | "OPTIONS";

export type HttpHeadersMap = { [key: string]: string };

export interface IHttpClient {
  req<T>(
    method: HttpMethod,
    path: string,
    body?: object,
    appendHeaders?: HttpHeadersMap,
  ): Promise<T>;
}

export interface HttpClientOptions {
  authorization?: string;
  headers?: HttpHeadersMap;
}

interface AccessToken extends AccessTokenModel {
  expiresDate: Date;
}

export class HttpClient implements IHttpClient {
  private accessToken: AccessToken | null = null;
  private accessTokenRequest: Promise<AccessTokenModel> | undefined;

  constructor(
    protected endpoint: string,
    private options = {} as HttpClientOptions,
  ) {}

  async req<T>(
    method: HttpMethod,
    path: string,
    body?: object,
    appendHeaders?: HttpHeadersMap,
  ): Promise<T> {
    const headers = new Headers();
    headers.set("Content-Type", "application/json");
    headers.set("Accept", "application/json");
    kv<string, string>(this.options?.headers).forEach(([k, v]) =>
      headers.set(k, v),
    );
    kv<string, string>(headers).forEach(([k, v]) => headers.set(k, v));

    if (this.options.authorization)
      headers.set("Authorization", this.options.authorization);
    else if (this.accessToken) {
      if (Date.now() - this.accessToken.expiresDate.getTime() > 0) {
        // Setting access token null here to avoid getting here again
        // right after we call getAndSetAccessToken().
        this.accessToken = null;
        return await this.getAndSetAccessToken(() =>
          this.req(method, path, body, appendHeaders),
        );
      }
      headers.set("Authorization", `accessToken ${this.accessToken.token}`);
    }
    const fullPath = replaceDoublePath(`${this.endpoint}/${path}`);
    const res = await window.fetch(fullPath, {
      method,
      headers,
      body: body ? JSON.stringify(body) : null,
      credentials: "include",
    });

    if (res.status === 204) {
      return {} as T;
    }

    let data = {};
    try {
      data = await res.json();
    } catch {
      /* empty */
    }

    if (
      res.status === 401 &&
      (data as ErrorResponse).error === "invalid access token"
    ) {
      return await this.getAndSetAccessToken(() =>
        this.req(method, path, body, appendHeaders),
      );
    }

    if (res.status >= 400) throw new APIError(res, data as ErrorResponse);

    return data as T;
  }

  private async getAccessToken(): Promise<AccessTokenModel> {
    if (!this.accessTokenRequest)
      this.accessTokenRequest = this.req("POST", "auth/accesstoken");
    return this.accessTokenRequest;
  }

  private async getAndSetAccessToken<T>(replay: () => Promise<T>): Promise<T> {
    const token = await this.getAccessToken();
    this.accessTokenRequest = undefined;
    this.accessToken = token as AccessToken;
    this.accessToken.expiresDate = new Date(token.expires);
    return await replay();
  }
}

function kv<TKey, TVal>(m?: any) {
  const _m = m ?? {};
  return Object.keys(_m).map((k) => [k, _m[k]]) as any as [TKey, TVal][];
}
