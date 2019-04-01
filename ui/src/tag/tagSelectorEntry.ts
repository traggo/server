import {Tags, Tags_tags} from '../gql/__generated__/Tags';
import {TagDefinitionType} from '../gql/__generated__/globalTypes';
import {QueryHookResult} from 'react-apollo-hooks';

export interface TagInputError {
    error: string;
    value: string;
}

export interface TagSelectorEntry {
    tag: Tags_tags & {create?: boolean; alreadyUsed?: boolean};
    value?: string;
}

export const specialTag = (name: string, state: 'used' | 'new'): TagSelectorEntry['tag'] => {
    return {
        key: name,
        __typename: 'TagDefinition',
        color: 'gray',
        type: TagDefinitionType.novalue,
        create: state === 'new',
        alreadyUsed: state === 'used',
    };
};

const tryAdd = (tagsResult: QueryHookResult<Tags, {}>, entry: string): TagSelectorEntry | TagInputError => {
    const [key, value, ...other] = entry.split(':');

    if (other.length || !tagsResult.data || !tagsResult.data.tags) {
        return {error: `${entry} has too many colons`, value: entry};
    }

    const foundTag = tagsResult.data.tags.find((tag) => tag.key === key);

    if (foundTag === undefined) {
        return {error: `'${key}' does not exist`, value: entry};
    }

    if (foundTag.type === TagDefinitionType.singlevalue && !value) {
        return {error: `'${key}' requires a value`, value: entry};
    }

    return {tag: foundTag, value};
};
interface EntriesAndErrors {
    errors: TagInputError[];
    entries: TagSelectorEntry[];
    usedTags: string[];
}

const groupAndCheckExistence = (a: EntriesAndErrors, entry: TagInputError | TagSelectorEntry) => {
    if ('tag' in entry) {
        if (a.usedTags.indexOf(entry.tag.key) === -1) {
            a.entries = [...a.entries, entry];
            a.usedTags = [...a.usedTags, entry.tag.key];
        } else {
            a.errors = [...a.errors, {value: label(entry), error: `'${entry.tag.key}' is already defined`}];
        }
    } else {
        a.errors = [...a.errors, entry];
    }
    return a;
};

export const addValues = (
    newValue: string,
    tagsResult: QueryHookResult<Tags, {}>,
    selectedEntries: TagSelectorEntry[]
): EntriesAndErrors => {
    return newValue
        .split(/\s+/)
        .filter((entry) => entry)
        .map((entry) => tryAdd(tagsResult, entry))
        .reduce(groupAndCheckExistence, {errors: [], entries: [], usedTags: selectedEntries.map((entry) => entry.tag.key)});
};

export const itemLabel = (tag: TagSelectorEntry) => {
    if (tag.tag.create) {
        return `Create tag '${tag.tag.key}'`;
    }
    if (tag.tag.alreadyUsed) {
        return `Tag '${tag.tag.key}' is already defined`;
    }
    return label(tag);
};

export const label = (tag: TagSelectorEntry) => {
    const suffix = tag.tag.type === TagDefinitionType.novalue ? '' : ':' + (tag.value || '');
    return `${tag.tag.key}${suffix}`;
};
