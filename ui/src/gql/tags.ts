import {gql} from 'apollo-boost';

export const SuggestTag = gql`
    query SuggestTag($query: String!) {
        tags: suggestTag(query: $query) {
            color
            key
        }
    }
`;

export const Tags = gql`
    query Tags {
        tags {
            color
            key
            usages
        }
    }
`;
export const SuggestTagValue = gql`
    query SuggestTagValue($tag: String!, $query: String!) {
        values: suggestTagValue(key: $tag, query: $query)
    }
`;

export const AddTag = gql`
    mutation AddTag($name: String!, $color: String!) {
        createTag(key: $name, color: $color) {
            color
            key
        }
    }
`;

export const UpdateTag = gql`
    mutation UpdateTag($key: String!, $newKey: String, $color: String!) {
        updateTag(key: $key, newKey: $newKey, color: $color) {
            color
            key
        }
    }
`;

export const RemoveTag = gql`
    mutation RemoveTag($key: String!) {
        removeTag(key: $key) {
            color
            key
        }
    }
`;
