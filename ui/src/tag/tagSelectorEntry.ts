import {Tags, Tags_tags} from '../gql/__generated__/Tags';
import {QueryResult} from 'react-apollo';

export interface TagInputError {
    error: string;
    value: string;
}

type SpecialTagState = 'used' | 'new' | 'no_values' | 'all_values_used';
export interface SpecialTag {
    create?: boolean;
    alreadyUsed?: boolean;
    noValues?: boolean;
    allValuesUsed?: boolean;
}

export interface TagSelectorEntry {
    tag: Omit<Tags_tags, 'usages'> & SpecialTag;
    value: string;
}

export interface InputTag {
    key: string;
    value: string;
}

export const toInputTags = (entries: TagSelectorEntry[]): InputTag[] => {
    return entries.map((entry) => ({key: entry.tag.key, value: entry.value}));
};

export const toTagSelectorEntry = (tags: Array<TagSelectorEntry['tag']>, entries: InputTag[]): TagSelectorEntry[] => {
    return entries.map(
        (timerTag): TagSelectorEntry => {
            const definition = tags.find((tag) => tag.key === timerTag.key) || specialTag(timerTag.key, 'new');
            return {
                tag: definition,
                value: timerTag.value,
            };
        }
    );
};

export const specialTag = (name: string, state: SpecialTagState): TagSelectorEntry['tag'] & {usages: 0} => {
    return {
        key: name,
        __typename: 'TagDefinition',
        color: 'gray',
        create: state === 'new',
        alreadyUsed: state === 'used',
        noValues: state === 'no_values',
        allValuesUsed: state === 'all_values_used',
        usages: 0,
    };
};

const tryAdd = (tagsResult: QueryResult<Tags, {}>, entry: string, onlyKeys: boolean): TagSelectorEntry | TagInputError => {
    const [keySomeCase, value, ...other] = entry.split(':');
    const key = keySomeCase.toLowerCase();

    if (other.length || !tagsResult.data || !tagsResult.data.tags) {
        return {error: `${entry} has too many colons`, value: entry};
    }

    const foundTag = tagsResult.data.tags.find((tag) => tag.key === key);

    if (foundTag === undefined) {
        return {error: `'${key}' does not exist`, value: entry};
    }

    if (onlyKeys && value) {
        return {error: `'${key}' has a value, but this field doesn't allow them`, value: entry};
    }

    if (!onlyKeys && !value) {
        return {error: `'${key}' requires a value`, value: entry};
    }

    return {tag: foundTag, value: value || ''};
};
interface EntriesAndErrors {
    errors: TagInputError[];
    entries: TagSelectorEntry[];
    usedTags: string[];
}

const groupAndCheckExistence = (onlyKeys: boolean, allowDuplicateKeys: boolean) => {
    return (a: EntriesAndErrors, entry: TagInputError | TagSelectorEntry) => {
        if ('tag' in entry) {
            if (allowDuplicateKeys || a.usedTags.indexOf(entry.tag.key) === -1) {
                a.entries = [...a.entries, entry];
                a.usedTags = [...a.usedTags, entry.tag.key];
            } else {
                a.errors = [...a.errors, {value: itemLabel(entry, onlyKeys), error: `'${entry.tag.key}' is already defined`}];
            }
        } else {
            a.errors = [...a.errors, entry];
        }
        return a;
    };
};

export const addValues = (
    newValue: string,
    tagsResult: QueryResult<Tags, {}>,
    selectedEntries: TagSelectorEntry[],
    onlyKeys: boolean,
    allowDuplicateKeys: boolean
): EntriesAndErrors => {
    return newValue
        .split(/\s+/)
        .filter((entry) => entry)
        .map((entry) => tryAdd(tagsResult, entry, onlyKeys))
        .reduce(groupAndCheckExistence(onlyKeys, allowDuplicateKeys), {
            errors: [],
            entries: [],
            usedTags: selectedEntries.map((entry) => entry.tag.key),
        });
};

export const itemLabel = (tag: TagSelectorEntry, onlyShowKey = false) => {
    if (tag.tag.create) {
        return `Create tag '${tag.tag.key}'`;
    }
    if (tag.tag.alreadyUsed) {
        return `Tag '${tag.tag.key}' is already defined`;
    }
    if (tag.tag.noValues) {
        return `Unkown value '${tag.value}' of tag '${tag.tag.key}'`;
    }
    if (tag.tag.allValuesUsed) {
        return `All values of tag '${tag.tag.key}' are used`;
    }
    if (onlyShowKey) {
        return tag.tag.key;
    }
    return label(tag);
};

export const label = (tag: TagSelectorEntry) => {
    return `${tag.tag.key}:${tag.value}`;
};
