import * as React from 'react';
import useTimeout from '@rooks/use-timeout';

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
