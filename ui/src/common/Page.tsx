import * as React from 'react';
import AppBar from '@material-ui/core/AppBar';
import CssBaseline from '@material-ui/core/CssBaseline';
import Divider from '@material-ui/core/Divider';
import Drawer from '@material-ui/core/Drawer';
import Hidden from '@material-ui/core/Hidden';
import IconButton from '@material-ui/core/IconButton';
import UsersIcon from '@material-ui/icons/SupervisorAccount';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import DashboardIcon from '@material-ui/icons/Dashboard';
import MenuIcon from '@material-ui/icons/Menu';
import SettingsIcon from '@material-ui/icons/Settings';
import DevicesIcon from '@material-ui/icons/DevicesOther';
import LabelIcon from '@material-ui/icons/Label';
import DashboardManageIcon from '@material-ui/icons/ListAlt';
import TimeLineIcon from '@material-ui/icons/Timeline';
import CalendarIcon from '@material-ui/icons/CalendarToday';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import {ListSubheader, Menu} from '@material-ui/core';
import HrefLink from '@material-ui/core/Link';
import * as gqlUser from '../gql/user';
import * as gqlVersion from '../gql/version';
import Button from '@material-ui/core/Button';
import AccountCircle from '@material-ui/icons/AccountCircle';
import {Link, matchPath, useLocation} from 'react-router-dom';
import MenuItem from '@material-ui/core/MenuItem';
import {useMutation, useQuery} from '@apollo/react-hooks';
import {Logout} from '../gql/__generated__/Logout';
import {Version} from '../gql/__generated__/Version';
import {CurrentUser} from '../gql/__generated__/CurrentUser';
import * as gqlDashboard from '../gql/dashboard';
import {Dashboards, Dashboards_dashboards} from '../gql/__generated__/Dashboards';
import makeStyles from '@material-ui/core/styles/makeStyles';

const drawerWidth = 240;

const useStyles = makeStyles((theme) => ({
    root: {
        display: 'flex',
        height: '100%',
    },
    drawer: {
        [theme.breakpoints.up('md')]: {
            width: drawerWidth,
            flexShrink: 0,
        },
    },
    appBar: {
        marginLeft: drawerWidth,
        [theme.breakpoints.up('md')]: {
            width: `calc(100% - ${drawerWidth}px)`,
        },
    },
    menuButton: {
        marginRight: 20,
        [theme.breakpoints.up('md')]: {
            display: 'none',
        },
    },
    toolbar: {
        ...theme.mixins.toolbar,
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'center',
    },
    drawerPaper: {
        width: drawerWidth,
    },
    content: {
        flexGrow: 1,
        paddingTop: theme.spacing(2),
        paddingBottom: theme.spacing(2),
        paddingLeft: theme.spacing(1),
        paddingRight: theme.spacing(1),
    },
    grow: {
        flexGrow: 1,
    },
    sectionDesktop: {
        display: 'none',
        [theme.breakpoints.up('xs')]: {
            display: 'flex',
        },
    },
    sectionMobile: {
        display: 'flex',
        [theme.breakpoints.up('xs')]: {
            display: 'none',
        },
    },
}));

// MUI v4's `component` prop isn't polymorphism-typed; revisited in the MUI v5 migration.
const routerLink = (to: string): any => {
    return React.forwardRef<HTMLAnchorElement>((props, ref) => <Link ref={ref} to={to} {...props} />);
};

const staticPageTitles: Array<{path: string; title: string}> = [
    {path: '/timesheet/list', title: 'Timesheet / List'},
    {path: '/timesheet/calendar', title: 'Timesheet / Calendar'},
    {path: '/user/settings', title: 'User / Settings'},
    {path: '/user/devices', title: 'User / Devices'},
    {path: '/user/tags', title: 'User / Tags'},
    {path: '/admin/users', title: 'Admin / Users'},
    {path: '/dashboards', title: 'Dashboards / Manage'},
];

const getPageTitle = (pathname: string, dashboards: Dashboards_dashboards[]): string => {
    for (const {path, title} of staticPageTitles) {
        if (matchPath(path, pathname)) {
            return title;
        }
    }
    const dashboardMatch = matchPath('/dashboard/:id/*', pathname);
    if (dashboardMatch && dashboardMatch.params.id) {
        const db = dashboards.find((dashboard) => dashboard.id === parseInt(dashboardMatch.params.id!, 10));
        return 'Dashboards / ' + (db ? db.name : '...');
    }
    return '';
};

export const Page: React.FC = ({children}) => {
    const classes = useStyles();
    const {data} = useQuery<CurrentUser>(gqlUser.CurrentUser);

    const [mobileOpen, setMobileOpen] = React.useState(false);
    const [userMenuOpen, setUserMenuOpen] = React.useState<null | HTMLElement>(null);
    const [logout] = useMutation<Logout>(gqlUser.Logout, {refetchQueries: [{query: gqlUser.CurrentUser}]});
    const {data: {version = gqlVersion.VersionDefault.version} = gqlVersion.VersionDefault} = useQuery<Version>(
        gqlVersion.Version
    );
    const dashboardsQuery = useQuery<Dashboards>(gqlDashboard.Dashboards);
    const dashboards = (dashboardsQuery.data && dashboardsQuery.data.dashboards) || [];

    const username = (data && data.user && data.user.name) || 'unknown';
    const admin = data && data.user && data.user.admin;

    const location = useLocation();
    const pageTitle = getPageTitle(location.pathname, dashboards);

    const drawer = (
        <div>
            <div className={classes.toolbar}>
                <HrefLink href="https://github.com/traggo" underline="none" target="_blank">
                    <Typography variant="h5" align="center" color="textPrimary">
                        traggo
                    </Typography>
                </HrefLink>
                <HrefLink href="https://github.com/traggo/server/releases" underline="none" target="_blank">
                    <Typography variant="subtitle2" align="center" color="textPrimary">
                        {version.name}@{version.commit.slice(0, 8)}
                    </Typography>
                </HrefLink>
            </div>
            <Divider />
            <List subheader={<ListSubheader>Dashboards</ListSubheader>} dense={true}>
                {dashboards.map(({id, name}) => (
                    <ListItem key={id} button component={routerLink(`/dashboard/${id}/${encodeURIComponent(name)}`)}>
                        <ListItemIcon>
                            <DashboardIcon />
                        </ListItemIcon>
                        <ListItemText primary={name} />
                    </ListItem>
                ))}
                {dashboards.length === 0 ? (
                    <ListItem button dense={true} disabled={true}>
                        <ListItemText primary={'no dashboards added'} />
                    </ListItem>
                ) : null}
                <ListItem button component={routerLink(`/dashboards`)}>
                    <ListItemIcon>
                        <DashboardManageIcon />
                    </ListItemIcon>
                    <ListItemText primary={'Manage'} />
                </ListItem>
            </List>
            <Divider />
            <List subheader={<ListSubheader>Timesheet</ListSubheader>} dense={true}>
                <ListItem button component={routerLink('/timesheet/list')}>
                    <ListItemIcon>
                        <TimeLineIcon />
                    </ListItemIcon>
                    <ListItemText primary="List" />
                </ListItem>
                <ListItem button component={routerLink('/timesheet/calendar')}>
                    <ListItemIcon>
                        <CalendarIcon />
                    </ListItemIcon>
                    <ListItemText primary="Calendar" />
                </ListItem>
            </List>
            <Divider />
            <List subheader={<ListSubheader>User</ListSubheader>} dense={true}>
                <ListItem button component={routerLink('/user/tags')}>
                    <ListItemIcon>
                        <LabelIcon />
                    </ListItemIcon>
                    <ListItemText primary="Tags" />
                </ListItem>
                <ListItem button component={routerLink('/user/devices')}>
                    <ListItemIcon>
                        <DevicesIcon />
                    </ListItemIcon>
                    <ListItemText primary="Devices" />
                </ListItem>
                <ListItem button component={routerLink('/user/settings')}>
                    <ListItemIcon>
                        <SettingsIcon />
                    </ListItemIcon>
                    <ListItemText primary="Settings" />
                </ListItem>
            </List>
            {admin ? (
                <>
                    <Divider />
                    <List subheader={<ListSubheader>Admin</ListSubheader>} dense={true}>
                        <ListItem button component={routerLink('/admin/users')}>
                            <ListItemIcon>
                                <UsersIcon />
                            </ListItemIcon>
                            <ListItemText primary="Users" />
                        </ListItem>
                    </List>
                </>
            ) : null}
            <Divider />
        </div>
    );

    return (
        <div className={classes.root}>
            <CssBaseline />
            <AppBar position="fixed" className={classes.appBar}>
                <Toolbar>
                    <IconButton
                        color="inherit"
                        aria-label="Open drawer"
                        onClick={() => setMobileOpen(!mobileOpen)}
                        className={classes.menuButton}>
                        <MenuIcon />
                    </IconButton>
                    <Typography variant="h6" color="inherit" noWrap>
                        {pageTitle}
                    </Typography>
                    <div className={classes.grow} />
                    <div className={classes.sectionDesktop}>
                        <Button color="inherit" onClick={(e) => setUserMenuOpen(e.currentTarget)}>
                            <AccountCircle />
                            &nbsp;{username}
                        </Button>
                        <Menu
                            anchorEl={userMenuOpen}
                            anchorOrigin={{vertical: 'top', horizontal: 'right'}}
                            transformOrigin={{vertical: 'top', horizontal: 'right'}}
                            open={!!userMenuOpen}
                            onClose={() => setUserMenuOpen(null)}>
                            <MenuItem onClick={() => setUserMenuOpen(null)} component={routerLink('/user/settings')}>
                                Settings
                            </MenuItem>
                            <MenuItem
                                onClick={() => {
                                    setUserMenuOpen(null);
                                    logout();
                                }}>
                                Logout
                            </MenuItem>
                        </Menu>
                    </div>
                </Toolbar>
            </AppBar>
            <nav className={classes.drawer}>
                <Hidden mdUp implementation="js">
                    <Drawer
                        variant="temporary"
                        open={mobileOpen}
                        onClose={() => setMobileOpen(false)}
                        classes={{
                            paper: classes.drawerPaper,
                        }}>
                        {drawer}
                    </Drawer>
                </Hidden>
                <Hidden smDown implementation="js">
                    <Drawer
                        classes={{
                            paper: classes.drawerPaper,
                        }}
                        variant="permanent"
                        open>
                        {drawer}
                    </Drawer>
                </Hidden>
            </nav>
            <main className={classes.content}>
                <div className={classes.toolbar} />
                {children}
            </main>
        </div>
    );
};
