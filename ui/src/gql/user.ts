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
        login(username: $name, pass: $pass, deviceName: "web ui", expiresAt: $expiresAt, cookie: true) {
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

export const Users = gql`
    query Users {
        users {
            id
            name
            admin
        }
        currentUser {
            id
        }
    }
`;

export const RemoveUser = gql`
    mutation RemoveUser($id: Int!) {
        removeUser(id: $id) {
            id
        }
    }
`;

export const UpdateUser = gql`
    mutation UpdateUser($id: Int!, $name: String!, $admin: Boolean!, $pass: String) {
        updateUser(id: $id, name: $name, admin: $admin, pass: $pass) {
            id
        }
    }
`;
export const CreateUser = gql`
    mutation CreateUser($name: String!, $admin: Boolean!, $pass: String!) {
        createUser(name: $name, admin: $admin, pass: $pass) {
            id
        }
    }
`;
