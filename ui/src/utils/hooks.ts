import * as React from 'react';
import useTimeout from '@rooks/use-timeout';

export const useStateAndDelegateWithDelayOnChange: <T>(
    initialState: T,
    delegate: (value: T) => void
) => [T | undefined, React.Dispatch<T>] = (initialState, delegate) => {
    const [value, setValue] = React.useState(initialState);
    const {start} = useTimeout(() => delegate(value), 50);
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
