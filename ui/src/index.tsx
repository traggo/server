import React from 'react';
import ReactDOM from 'react-dom';
import {Root} from './Root';
import moment from 'moment-timezone';

moment.updateLocale('en', {
    week: {
        dow: 1,
        doy: moment.localeData('en').firstDayOfYear(),
    },
});

ReactDOM.render(<Root />, document.getElementById('root'));
