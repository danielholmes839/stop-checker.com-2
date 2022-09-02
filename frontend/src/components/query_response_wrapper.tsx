import { UseQueryState } from "urql";

export const QueryResponseWrapper: React.FC<{ response: UseQueryState }> = ({
  response,
  children,
}) => {
  const { data, fetching, error } = response;

  if (fetching) {
    return <>Loading...</>;
  }

  if (error || !data) {
    return <>error</>;
  }

  return <>{children}</>;
};
