import { Container } from "components";
import { Search, StopPreviewDefaultActions } from "./search";

export { Search } from "./search";

export const SearchPage = () => {
  return (
    <Container>
      <div className="mt-3">
        <h1 className="text-3xl font-semibold mb-3">Browse Stops</h1>
        <Search
          config={{
            Actions: StopPreviewDefaultActions,
            enableMap: true,
            enableStopRouteLinks: true,
          }}
        />
      </div>
    </Container>
  );
};
