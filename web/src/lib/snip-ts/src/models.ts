/** @format */

export interface AccessTokenModel {
  token: string;
  expires: string;
}

export interface CodeResponse {
  code: number;
}

export interface ErrorResponse extends CodeResponse {
  error: string;
}
