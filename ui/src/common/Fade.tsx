import * as React from 'react';
import {CSSTransition} from 'react-transition-group';

const duration = 100;

const defaultStyle = {
    transition: `all ${duration}ms ease-in-out`,
    opacity: 0,
};

export const Fade: React.FC<{fullyVisible: boolean; opacity?: number}> = ({fullyVisible, children, opacity = 0.4}) => {
    const transitionStyles: Record<string, React.CSSProperties> = {
        entering: {opacity: 1},
        entered: {opacity: 1},
        exiting: {opacity},
        exited: {opacity},
    };
    return (
        <CSSTransition in={fullyVisible} timeout={duration}>
            {(state) => (
                <div
                    style={{
                        ...defaultStyle,
                        ...transitionStyles[state],
                    }}>
                    {children}
                </div>
            )}
        </CSSTransition>
    );
};
