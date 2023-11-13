import { AccessTokenModel, CodeResponse } from "./models";

import { Client } from "./client";

import { SubClient } from "./subclient";

export class AuthClient extends SubClient {
  constructor(client: Client) {
    super(client, "auth");
  }

  accesstoken(): Promise<AccessTokenModel> {
    return this.req("POST", "accesstoken");
  }

  logout(): Promise<CodeResponse> {
    return this.req("POST", "logout");
  }

  signup(): Promise<AccessTokenModel> {
    return this.req("POST", "signup");
  }
}
