import {gql} from 'apollo-boost';

export const SuggestTag = gql`
    query SuggestTag($query: String!) {
        tags: suggestTag(query: $query) {
            color
            key
            type
        }
    }
`;

export const Tags = gql`
    query Tags {
        tags {
            color
            key
            type
        }
    }
`;
export const SuggestTagValue = gql`
    query SuggestTagValue($tag: String!, $query: String!) {
        values: suggestTagValue(key: $tag, query: $query)
    }
`;

export const AddTag = gql`
    mutation AddTag($name: String!, $color: String!, $type: TagDefinitionType!) {
        createTag(key: $name, color: $color, type: $type) {
            color
            key
            type
        }
    }
`;
