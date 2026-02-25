export const customInstance = async <T>(
  url: string,
  options?: RequestInit,
): Promise<T> => {
  const response = await fetch(`/api/public/v1${url}`, {
    credentials: "include",
    ...options,
  });

  if (!response.ok) {
    throw response;
  }

  return response.json() as Promise<T>;
};

export type ErrorType<Error> = Error;
export type BodyType<BodyData> = BodyData;
