import * as React from 'react';
import AppBar from '@mui/material/AppBar';
import CssBaseline from '@mui/material/CssBaseline';
import Divider from '@mui/material/Divider';
import Drawer from '@mui/material/Drawer';
import Hidden from '@mui/material/Hidden';
import IconButton from '@mui/material/IconButton';
import UsersIcon from '@mui/icons-material/SupervisorAccount';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import DashboardIcon from '@mui/icons-material/Dashboard';
import MenuIcon from '@mui/icons-material/Menu';
import SettingsIcon from '@mui/icons-material/Settings';
import DevicesIcon from '@mui/icons-material/DevicesOther';
import LabelIcon from '@mui/icons-material/Label';
import DashboardManageIcon from '@mui/icons-material/ListAlt';
import TimeLineIcon from '@mui/icons-material/Timeline';
import CalendarIcon from '@mui/icons-material/CalendarToday';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';
import {ListSubheader, Menu} from '@mui/material';
import HrefLink from '@mui/material/Link';
import * as gqlUser from '../gql/user';
import * as gqlVersion from '../gql/version';
import Button from '@mui/material/Button';
import AccountCircle from '@mui/icons-material/AccountCircle';
import {Link, matchPath, useLocation} from 'react-router-dom';
import MenuItem from '@mui/material/MenuItem';
import {useMutation, useQuery} from '@apollo/client';
import {LogoutMutation, VersionQuery, CurrentUserQuery, DashboardsQuery} from '../gql/__generated__';
import * as gqlDashboard from '../gql/dashboard';
import {Dashboard} from '../gql/types';
import makeStyles from '@mui/styles/makeStyles';

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

// MUI's `component` prop isn't polymorphism-typed against an arbitrary forwardRef component.
// eslint-disable-next-line @typescript-eslint/no-explicit-any
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

const getPageTitle = (pathname: string, dashboards: Dashboard[]): string => {
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

export const Page: React.FC<React.PropsWithChildren> = ({children}) => {
    const classes = useStyles({});
    const {data} = useQuery<CurrentUserQuery>(gqlUser.CurrentUser);

    const [mobileOpen, setMobileOpen] = React.useState(false);
    const [userMenuOpen, setUserMenuOpen] = React.useState<null | HTMLElement>(null);
    const [logout] = useMutation<LogoutMutation>(gqlUser.Logout, {refetchQueries: [{query: gqlUser.CurrentUser}]});
    const {data: {version = gqlVersion.VersionDefault.version} = gqlVersion.VersionDefault} = useQuery<VersionQuery>(
        gqlVersion.Version
    );
    const dashboardsQuery = useQuery<DashboardsQuery>(gqlDashboard.Dashboards);
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
                        className={classes.menuButton}
                        size="large">
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
                <Hidden mdDown implementation="js">
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
