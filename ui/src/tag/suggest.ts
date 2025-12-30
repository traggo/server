import {Tags} from '../gql/__generated__/Tags';
import * as gqlTags from '../gql/tags';
import {useQuery} from '@apollo/react-hooks';
import {SuggestTagValue, SuggestTagValueVariables} from '../gql/__generated__/SuggestTagValue';
import {TagSelectorEntry, specialTag} from './tagSelectorEntry';
import {QueryResult} from 'react-apollo';

export const useSuggest = (
    tagResult: QueryResult<Tags, {}>,
    inputValue: string,
    usedTags: TagSelectorEntry[],
    skipValue = false,
    allowDuplicateKeys = false,
    createTags = true
): TagSelectorEntry[] => {
    const [tagKeySomeCase, tagValue] = inputValue.split(':');
    const tagKey = tagKeySomeCase.toLowerCase();

    const exactMatch = ((tagResult.data && tagResult.data.tags) || []).find((tag) => tag.key === tagKey);

    const valueResult = useQuery<SuggestTagValue, SuggestTagValueVariables>(gqlTags.SuggestTagValue, {
        variables: {tag: tagKey, query: tagValue},
        skip: exactMatch === undefined || skipValue,
    });

    const usedKeys = usedTags.map((t) => t.tag.key);
    // Filter usedValues to only include values from the same tag key
    const usedValues = usedTags.filter((t) => t.tag.key === tagKey).map((t) => t.value);

    if (exactMatch && tagValue !== undefined && !skipValue && (allowDuplicateKeys || usedKeys.indexOf(exactMatch.key) === -1)) {
        return suggestTagValue(exactMatch, tagValue, valueResult, usedValues, createTags);
    } else {
        return suggestTag(exactMatch, tagResult, tagKey, usedKeys, allowDuplicateKeys, createTags);
    }
};

const suggestTag = (
    exactMatch: TagSelectorEntry['tag'] | undefined,
    tagResult: QueryResult<Tags, {}>,
    tagKey: string,
    usedTags: string[],
    allowDuplicateKeys: boolean,
    createTags: boolean
) => {
    if (!tagResult.data || tagResult.data.tags === null) {
        return [];
    }

    let availableTags = (tagResult.data.tags || [])
        .filter((tag) => allowDuplicateKeys || usedTags.indexOf(tag.key) === -1)
        .filter((tag) => tag.key.indexOf(tagKey) === 0);

    if (tagKey && !exactMatch && createTags) {
        availableTags = [specialTag(tagKey, 'new'), ...availableTags];
    }

    if (usedTags.indexOf(tagKey) !== -1 && !allowDuplicateKeys) {
        availableTags = [specialTag(tagKey, 'used'), ...availableTags];
    }

    return availableTags
        .sort((a, b) => b.usages - a.usages)
        .slice(0, 5)
        .map((tag) => ({tag, value: ''}));
};

const suggestTagValue = (
    exactMatch: TagSelectorEntry['tag'],
    tagValue: string,
    valueResult: QueryResult<SuggestTagValue, SuggestTagValueVariables>,
    usedValues: string[],
    includeInputValueOnNoMatch: boolean
): TagSelectorEntry[] => {
    let someValues = (valueResult.data && valueResult.data.values) || [];

    if (includeInputValueOnNoMatch && someValues.indexOf(tagValue) === -1) {
        someValues = [tagValue, ...someValues];
    }

    if (someValues.length === 0 && !includeInputValueOnNoMatch) {
        return [{tag: specialTag(exactMatch.key, 'no_values'), value: tagValue}];
    }

    someValues = someValues.filter((val) => usedValues.indexOf(val) === -1);
    if (someValues.length === 0 && !includeInputValueOnNoMatch) {
        return [{tag: specialTag(exactMatch.key, 'all_values_used'), value: ''}];
    }

    return someValues.map((val) => ({tag: exactMatch, value: val}));
};
