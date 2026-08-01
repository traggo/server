import moment from 'moment';

// Registers the locales via moment.defineLocale() directly (config copied verbatim from
// moment's own locale/de.js, locale/en-au.js and locale/en-gb.js) instead of importing
// 'moment/locale/de' etc. Those deep imports are UMD files with an internal
// `require('../moment')` - under Vite/Rollup's commonjs interop that relative require can end
// up resolving to a *different* bundled instance of moment than the rest of the app's plain
// `import moment from 'moment'`, so the locale silently registers onto an instance nothing else
// ever reads from. Calling defineLocale ourselves, through the same import every other module
// uses, sidesteps that entirely.

function processGermanRelativeTime(number: number, withoutSuffix: boolean, key: string): string {
    const format: Record<string, [string, string]> = {
        m: ['eine Minute', 'einer Minute'],
        h: ['eine Stunde', 'einer Stunde'],
        d: ['ein Tag', 'einem Tag'],
        dd: [number + ' Tage', number + ' Tagen'],
        w: ['eine Woche', 'einer Woche'],
        M: ['ein Monat', 'einem Monat'],
        MM: [number + ' Monate', number + ' Monaten'],
        y: ['ein Jahr', 'einem Jahr'],
        yy: [number + ' Jahre', number + ' Jahren'],
    };
    return withoutSuffix ? format[key][0] : format[key][1];
}

const ordinalEn = (number: number): string => {
    if (~~((number % 100) / 10) === 1) {
        return number + 'th';
    }
    switch (number % 10) {
        case 1:
            return number + 'st';
        case 2:
            return number + 'nd';
        case 3:
            return number + 'rd';
        default:
            return number + 'th';
    }
};

moment.defineLocale('de', {
    months: 'Januar_Februar_März_April_Mai_Juni_Juli_August_September_Oktober_November_Dezember'.split('_'),
    monthsShort: 'Jan._Feb._März_Apr._Mai_Juni_Juli_Aug._Sep._Okt._Nov._Dez.'.split('_'),
    monthsParseExact: true,
    weekdays: 'Sonntag_Montag_Dienstag_Mittwoch_Donnerstag_Freitag_Samstag'.split('_'),
    weekdaysShort: 'So._Mo._Di._Mi._Do._Fr._Sa.'.split('_'),
    weekdaysMin: 'So_Mo_Di_Mi_Do_Fr_Sa'.split('_'),
    weekdaysParseExact: true,
    longDateFormat: {
        LT: 'HH:mm',
        LTS: 'HH:mm:ss',
        L: 'DD.MM.YYYY',
        LL: 'D. MMMM YYYY',
        LLL: 'D. MMMM YYYY HH:mm',
        LLLL: 'dddd, D. MMMM YYYY HH:mm',
    },
    calendar: {
        sameDay: '[heute um] LT [Uhr]',
        sameElse: 'L',
        nextDay: '[morgen um] LT [Uhr]',
        nextWeek: 'dddd [um] LT [Uhr]',
        lastDay: '[gestern um] LT [Uhr]',
        lastWeek: '[letzten] dddd [um] LT [Uhr]',
    },
    relativeTime: {
        future: 'in %s',
        past: 'vor %s',
        s: 'ein paar Sekunden',
        ss: '%d Sekunden',
        m: processGermanRelativeTime,
        mm: '%d Minuten',
        h: processGermanRelativeTime,
        hh: '%d Stunden',
        d: processGermanRelativeTime,
        dd: processGermanRelativeTime,
        w: processGermanRelativeTime,
        ww: '%d Wochen',
        M: processGermanRelativeTime,
        MM: processGermanRelativeTime,
        y: processGermanRelativeTime,
        yy: processGermanRelativeTime,
    },
    dayOfMonthOrdinalParse: /\d{1,2}\./,
    ordinal: (n: number) => n + '.',
    week: {
        dow: 1,
        doy: 4,
    },
});

moment.defineLocale('en-au', {
    months: 'January_February_March_April_May_June_July_August_September_October_November_December'.split('_'),
    monthsShort: 'Jan_Feb_Mar_Apr_May_Jun_Jul_Aug_Sep_Oct_Nov_Dec'.split('_'),
    weekdays: 'Sunday_Monday_Tuesday_Wednesday_Thursday_Friday_Saturday'.split('_'),
    weekdaysShort: 'Sun_Mon_Tue_Wed_Thu_Fri_Sat'.split('_'),
    weekdaysMin: 'Su_Mo_Tu_We_Th_Fr_Sa'.split('_'),
    longDateFormat: {
        LT: 'h:mm A',
        LTS: 'h:mm:ss A',
        L: 'DD/MM/YYYY',
        LL: 'D MMMM YYYY',
        LLL: 'D MMMM YYYY h:mm A',
        LLLL: 'dddd, D MMMM YYYY h:mm A',
    },
    calendar: {
        sameDay: '[Today at] LT',
        nextDay: '[Tomorrow at] LT',
        nextWeek: 'dddd [at] LT',
        lastDay: '[Yesterday at] LT',
        lastWeek: '[Last] dddd [at] LT',
        sameElse: 'L',
    },
    relativeTime: {
        future: 'in %s',
        past: '%s ago',
        s: 'a few seconds',
        ss: '%d seconds',
        m: 'a minute',
        mm: '%d minutes',
        h: 'an hour',
        hh: '%d hours',
        d: 'a day',
        dd: '%d days',
        M: 'a month',
        MM: '%d months',
        y: 'a year',
        yy: '%d years',
    },
    dayOfMonthOrdinalParse: /\d{1,2}(st|nd|rd|th)/,
    ordinal: ordinalEn,
    week: {
        dow: 0,
        doy: 4,
    },
});

moment.defineLocale('en-gb', {
    months: 'January_February_March_April_May_June_July_August_September_October_November_December'.split('_'),
    monthsShort: 'Jan_Feb_Mar_Apr_May_Jun_Jul_Aug_Sep_Oct_Nov_Dec'.split('_'),
    weekdays: 'Sunday_Monday_Tuesday_Wednesday_Thursday_Friday_Saturday'.split('_'),
    weekdaysShort: 'Sun_Mon_Tue_Wed_Thu_Fri_Sat'.split('_'),
    weekdaysMin: 'Su_Mo_Tu_We_Th_Fr_Sa'.split('_'),
    longDateFormat: {
        LT: 'HH:mm',
        LTS: 'HH:mm:ss',
        L: 'DD/MM/YYYY',
        LL: 'D MMMM YYYY',
        LLL: 'D MMMM YYYY HH:mm',
        LLLL: 'dddd, D MMMM YYYY HH:mm',
    },
    calendar: {
        sameDay: '[Today at] LT',
        nextDay: '[Tomorrow at] LT',
        nextWeek: 'dddd [at] LT',
        lastDay: '[Yesterday at] LT',
        lastWeek: '[Last] dddd [at] LT',
        sameElse: 'L',
    },
    relativeTime: {
        future: 'in %s',
        past: '%s ago',
        s: 'a few seconds',
        ss: '%d seconds',
        m: 'a minute',
        mm: '%d minutes',
        h: 'an hour',
        hh: '%d hours',
        d: 'a day',
        dd: '%d days',
        M: 'a month',
        MM: '%d months',
        y: 'a year',
        yy: '%d years',
    },
    dayOfMonthOrdinalParse: /\d{1,2}(st|nd|rd|th)/,
    ordinal: ordinalEn,
    week: {
        dow: 1,
        doy: 4,
    },
});
