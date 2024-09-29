import * as React from 'react';
import {makeStyles} from '@material-ui/core';

const useStyle = makeStyles((theme) => {
    return {
        root: {
            height: '100%',
            width: '100%',
            // '& .fc-time-grid-event': {
            //     opacity: 0,
            //     border: 'none',
            // },
            '& .fc-head': {
                fontFamily: 'monospace',
            },
            '& .fc-toolbar': {
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
            },
            '& .fc-left, .fc-center, .fc-right': {
                display: 'inline-block',
                verticalAlign: 'middle',
                textAlign: 'center',
                // width: '33.33%',
                boxSizing: 'border-box',
                // padding: '0 10px',
            },
            '& .fc-center h2': {
                margin: 0,
                fontSize: '1.5rem',
            },
            '@media (max-width: 600px)': {
                '& .fc-toolbar': {
                    display: 'inherit',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                },
                '& .fc-left, .fc-center, .fc-right': {
                    display: 'block',
                    width: '100%',
                },
                '& .fc-left, .fc-right': {
                    textAlign: 'center',
                },
            },
            '& .fc-time-grid-event.fc-v-event.fc-event': {
                borderRadius: 1,
                border: 'none',
                padding: 5,
                opacity: 0.65,
                left: '5% !important',
                right: '5% !important',
            },
            '& .fc .fc-time-grid-event': {
                minHeight: 1,
            },
            '& .fc-event': {
                fontSize: 'inherit !important',
                borderRadius: 1,
            },
            '& .fc td, .fc th': {
                borderStyle: 'solid !important',
                borderWidth: '0px !important',
                padding: '0 !important',
                verticalAlign: 'top !important',
            },
            '& .fc tr:nth-child(4n)': {
                borderStyle: 'solid !important',
                borderColor: theme.palette.type === 'dark' ? 'rgba(255,255,255,.2)' : 'rgba(0,0,0,.2)',
                borderWidth: '0 0 2px 0 !important',
                padding: '0 !important',
                verticalAlign: 'top !important',
            },
            '& .fc tbody tr:first-child': {
                borderStyle: 'solid !important',
                borderColor: theme.palette.type === 'dark' ? 'rgba(255,255,255,.2)' : 'rgba(0,0,0,.2)',
                borderWidth: '2px 0 0 0 !important',
                padding: '0 !important',
                verticalAlign: 'top !important',
            },
            '& .fc td': {
                borderStyle: 'solid !important',
                borderColor: theme.palette.type === 'dark' ? 'rgba(255,255,255,.1)' : 'rgba(0,0,0,.1)',
                borderWidth: '0 1px 0 0 !important',
                padding: '0 !important',
                verticalAlign: 'top !important',
            },
            '& .fc tr:nth-child(4n+2)': {
                borderStyle: 'solid !important',
                borderColor: theme.palette.type === 'dark' ? 'rgba(255,255,255,.1)' : 'rgba(0,0,0,.1)',
                borderWidth: '0 0 1px 0 !important',
                padding: '0 !important',
                verticalAlign: 'top !important',
            },
            '& .fc-event .fc-bg': {
                zIndex: '1 !important',
                background: 'inherit !important',
                opacity: '.25 !important',
            },
            '& .fc-time-grid-event .fc-time': {
                fontWeight: 'normal !important',
            },
            '& .fc-ltr .fc-h-event.fc-not-end, .fc-rtl .fc-h-event.fc-not-start': {
                opacity: '.65 !important',
                marginLeft: '12px !important',
                padding: '5px !important',
            },
            '& .fc-day-grid-event.fc-h-event.fc-event.fc-not-start.fc-end': {
                opacity: '.65 !important',
                marginLeft: '12px !important',
                padding: '5px !important',
            },
            '& .fc-today': {
                background: theme.palette.type === 'dark' ? 'rgba(199,199,199,.1) !important' : 'rgba(0,0,0,.1) !important',
            },
            '& .fc-button': {
                display: 'inline-block',
                position: 'relative',
                cursor: 'pointer',
                minHeight: '36px',
                minWidth: '88px',
                lineHeight: '36px',
                verticalAlign: 'middle',
                alignItems: 'center',
                textAlign: 'center',
                borderRadius: '2px',
                boxSizing: 'border-box',
                userSelect: 'none',
                outline: 'none',
                border: '0',
                padding: '0 6px',
                margin: '6px 8px',
                letterSpacing: '.01em',
                background: 'transparent',
                color: 'currentColor',
                whiteSpace: 'nowrap',
                textTransform: 'uppercase',
                fontWeight: '500',
                fontSize: '14px',
                fontStyle: 'inherit',
                fontVariant: 'inherit',
                fontFamily: 'inherit',
                textDecoration: 'none',
                overflow: 'hidden',
                borderColor: theme.palette.primary.main,
                borderWidth: 1,
                borderStyle: 'solid',
                transition: 'box-shadow .4s cubic-bezier(.25,.8,.25,1),background-color .4s cubic-bezier(.25,.8,.25,1)',
            },
            '& .fc-button-active': {
                backgroundColor: theme.palette.primary.main + '!important',
            },
            '& .fc-button:active': {
                backgroundColor: theme.palette.primary.main + '!important',
            },
            '& .fc td.fc-axis.fc-time': {
                position: 'absolute',
                marginTop: -10,
                background: theme.palette.background.paper,
                zIndex: 1000,
            },
            '& .fc td.fc-axis.fc-time span': {
                paddingRight: 10,
            },
            '& .fc .fc-time-grid-event-inset': {
                boxShadow: theme.shadows[5],
            },
            '& .fc .fc-mirror': {
                opacity: 0.5,
            },
            '& .fc .fc-time-grid-event:hover': {
                minHeight: 17,
                zIndex: '2!important',
            },
            '& .fc .fc-time-grid-event[data-has-end="false"]': {
                minHeight: 40,
            },
            '& .fc .fc-content': {
                height: '100%',
            },
            '& .fc .fc-content .stop': {
                background: theme.palette.type === 'dark' ? 'rgba(255,255,255,.2) !important' : 'rgba(0,0,0,.1) !important',
                borderRadius: '6px 0 0 0',
                fontSize: '10px',
                width: '50%',
                padding: '2px 0 1px 0',
                textAlign: 'center',
                right: 0,
                bottom: 0,
                position: 'absolute',
                display: 'inline-block',
            },
            '& .fc .fc-content .stop:hover': {
                background: theme.palette.type === 'dark' ? 'rgba(255,255,255,.3) !important' : 'rgba(0,0,0,.2) !important',
            },
            '& .fc .fc-content .stop a': {
                width: '100%',
                display: 'block',
            },
            '& .fc .fc-now-indicator': {
                zIndex: 5,
            },
            '& .fc .ellipsis': {
                display: '-webkit-box',
                boxOrient: 'vertical',
                overflow: 'hidden',
            },
            '& .fc .ellipsis-single': {
                textOverflow: 'ellipsis',
                display: 'inline-block',
                width: '100%',
                whiteSpace: 'nowrap',
                overflow: 'hidden',
            },
            '& .fc .__start': {
                background: theme.palette.primary.main,
                color: theme.palette.primary.contrastText,
                border: 0,
                borderRadius: '0 0 10px 10px',
                textAlign: 'center',
                boxShadow: theme.shadows[1],
                cursor: 'pointer',
            },
            '& .fc .__start:hover': {
                background: theme.palette.primary.light,
            },
            // tslint:disable-next-line:no-any important breaks it
        } as any,
    };
});

export const FullCalendarStyling: React.FC = ({children}) => {
    const classes = useStyle();
    return <div className={classes.root}>{children}</div>;
};
