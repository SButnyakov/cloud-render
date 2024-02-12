export type RefreshTokenResponse = {
  'refresh_token': string,
  'access_token': string
}

export enum SigninResponseCodes {
  INVALID_CREDENTIALS = 400,
  INTERNAL_SERVER_ERROR = 501
}

export enum SignupResponseCodes {
  INTERNAL_SERVER_ERROR = 501
}