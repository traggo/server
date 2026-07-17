import * as React from 'react';

interface UseTimeoutHandler {
    start: () => void;
    stop: () => void;
    isActive: boolean;
}

export const useTimeout = (callback: () => void, delayMs?: number): UseTimeoutHandler => {
    const savedCallback = React.useRef(callback);
    savedCallback.current = callback;
    const handle = React.useRef<number | undefined>(undefined);
    const [isActive, setIsActive] = React.useState(false);

    const stop = React.useCallback(() => {
        window.clearTimeout(handle.current);
        setIsActive(false);
    }, []);

    const start = React.useCallback(() => {
        window.clearTimeout(handle.current);
        setIsActive(true);
        handle.current = window.setTimeout(() => {
            setIsActive(false);
            savedCallback.current();
        }, delayMs);
    }, [delayMs]);

    React.useEffect(() => () => window.clearTimeout(handle.current), []);

    return {start, stop, isActive};
};

export const useInterval = (callback: () => void, intervalDurationMs: number, startImmediate = false): void => {
    const savedCallback = React.useRef(callback);
    savedCallback.current = callback;

    React.useEffect(() => {
        if (startImmediate) {
            savedCallback.current();
        }
        const handle = window.setInterval(() => savedCallback.current(), intervalDurationMs);
        return () => window.clearInterval(handle);
    }, [intervalDurationMs, startImmediate]);
};

export const useStateAndDelegateWithDelayOnChange: <T>(
    initialState: T,
    delegate: React.Dispatch<React.SetStateAction<T>>,
    delay?: number
) => [T, React.Dispatch<React.SetStateAction<T>>] = (initialState, delegate, delay = 50) => {
    const [value, setValue] = React.useState(initialState);
    const {start} = useTimeout(() => delegate(value), delay);
    return [
        value,
        (newValue) => {
            if (value !== newValue) {
                start();
            }
            setValue(newValue);
        },
    ];
};
