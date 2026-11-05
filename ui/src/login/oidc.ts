import * as React from 'react';

export interface OidcInfo {
    enabled: boolean;
    passwordLogin: boolean;
    loginUrl?: string;
}

// OIDC is optional, if the info endpoint is not reachable the regular password login is used.
const fallback: OidcInfo = {enabled: false, passwordLogin: true};

export const useOidcInfo = (): OidcInfo | undefined => {
    const [info, setInfo] = React.useState<OidcInfo | undefined>(undefined);
    React.useEffect(() => {
        let active = true;
        fetch('/auth/oidc/info', {credentials: 'same-origin'})
            .then((response) => (response.ok ? response.json() : fallback))
            .catch(() => fallback)
            .then((result: OidcInfo) => active && setInfo(result));
        return () => {
            active = false;
        };
    }, []);
    return info;
};
