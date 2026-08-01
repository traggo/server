import * as React from 'react';

export const Center: React.FC<React.PropsWithChildren> = ({children}) => {
    return <div style={{display: 'flex', alignItems: 'center', justifyContent: 'center', height: '100%'}}>{children}</div>;
};
