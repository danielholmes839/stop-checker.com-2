import { StopPreviewFragment } from "client/types";
import { Card } from "components";
import { Search } from "pages/search";
import { StopPreviewActions } from "pages/search/search";
import React, { useState } from "react";
import { Wizard, useWizard } from "react-use-wizard";

const Actions: StopPreviewActions = ({ stop }) => {
  const { origin, setOrigin, destination, setDestination } = useAutomatic();

  const selectedAsOrigin = origin !== null && stop.id === origin.id;
  const selectedAsDestination =
    destination !== null && stop.id === destination.id;

  return (
    <div className="mt-3">
      <>
        <button
          disabled={selectedAsOrigin}
          className={
            selectedAsOrigin
              ? "mr-3 text-gray-500 text-sm"
              : "mr-3 text-primary-500 underline text-sm"
          }
          onClick={() => setOrigin(stop)}
        >
          {selectedAsOrigin ? "Selected as Origin" : "Set as Origin"}
        </button>
        <button
          className={
            selectedAsDestination
              ? "mr-3 text-gray-500 text-sm"
              : "mr-3 text-primary-500 underline text-sm"
          }
          onClick={() => setDestination(stop)}
        >
          {selectedAsDestination
            ? "Selected as Destination"
            : "Set as Destination"}
        </button>
      </>
    </div>
  );
};

type AutomaticContextValue = {
  origin: StopPreviewFragment | null;
  setOrigin: (stop: StopPreviewFragment) => void;
  destination: StopPreviewFragment | null;
  setDestination: (stop: StopPreviewFragment) => void;
  swap: () => void;
};

const AutomaticContext = React.createContext<AutomaticContextValue>({
  origin: null,
  setOrigin: (stop) => {},
  destination: null,
  setDestination: (stop) => {},
  swap: () => {},
});

const useAutomatic = () => React.useContext(AutomaticContext);

export const AutomaticProvider: React.FC = ({ children }) => {
  const [origin, setOrigin] = useState<StopPreviewFragment | null>(null);
  const [destination, setDestination] = useState<StopPreviewFragment | null>(
    null
  );
  const swap = () => {
    let originCopy = origin ? { ...origin } : null;
    setOrigin(destination);
    setDestination(originCopy);
  };

  return (
    <AutomaticContext.Provider
      value={{
        origin: origin,
        setOrigin: (stop) => {
          if (destination && destination.id === stop.id) {
            swap();
            return;
          }
          setOrigin(stop);
        },
        destination: destination,
        setDestination: (stop) => {
          if (origin && origin.id === stop.id) {
            swap();
            return;
          }
          setDestination(stop);
        },
        swap: swap,
      }}
    >
      {children}
    </AutomaticContext.Provider>
  );
};

const Current: React.FC = () => {
  const { origin, destination } = useAutomatic();
  const { nextStep } = useWizard();
  return (
    <div className="mt-3">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-x-5 gap-y-3">
        <Card>
          <h1>
            <div className="bg-red-600 rounded-full px-1 text-white text-xs inline mr-2 font-bold">
              A
            </div>
            {origin ? (
              <span>{origin.name} (Origin)</span>
            ) : (
              <span>Origin not selected</span>
            )}
          </h1>
        </Card>
        <Card>
          <h1>
            <div className="bg-red-600 rounded-full px-1 text-white text-xs inline mr-2 font-bold">
              B
            </div>{" "}
            {destination ? (
              <span>{destination.name} (Destination)</span>
            ) : (
              <span>Destination not selected</span>
            )}
          </h1>
        </Card>
      </div>
      <button className="mt-1 bg-gray-100" onClick={nextStep}>
        Next
      </button>
    </div>
  );
};

const Setup: React.FC = () => {
  return (
    <div>
      <Current />
      <div className="mt-3">
        <Search
          config={{
            Actions: Actions,
            enableMap: true,
            enableStopRouteLinks: false,
          }}
        />
      </div>
    </div>
  );
};

const Results: React.FC = () => {
  return <>Results</>;
};

export const Automatic: React.FC = () => {
  return (
    <AutomaticProvider>
      <div className="mt-3">
        <h1 className="text-3xl font-semibold">Travel Planner</h1>
      </div>
      <Wizard>
        <Setup />
        <Results />
      </Wizard>
    </AutomaticProvider>
  );
};