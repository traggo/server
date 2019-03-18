import {gql} from 'apollo-boost';

export const CurrentUser = gql`
    query CurrentUser {
        user: currentUser {
            name
            admin
            id
        }
    }
`;

export const Login = gql`
    mutation Login($name: String!, $pass: String!, $expiresAt: Time!) {
        login: createDevice(username: $name, pass: $pass, deviceName: "web ui", expiresAt: $expiresAt, cookie: true) {
            user {
                id
                name
                admin
            }
        }
    }
`;
export const Logout = gql`
    mutation Logout {
        user: removeCurrentDevice {
            name
        }
    }
`;
