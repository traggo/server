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
import ViewLintIcon from '@material-ui/icons/ViewList';
import CalendarIcon from '@material-ui/icons/CalendarToday';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import {StyleRulesCallback, WithStyles, withStyles} from '@material-ui/core/styles';
import {ListSubheader, Menu} from '@material-ui/core';
import HrefLink from '@material-ui/core/Link';
import * as gqlUser from '../gql/user';

import Button from '@material-ui/core/Button';
import AccountCircle from '@material-ui/icons/AccountCircle';
import HighlightIcon from '@material-ui/icons/Highlight';
import {Link} from 'react-router-dom';
import MenuItem from '@material-ui/core/MenuItem';
import {useMutation} from 'react-apollo-hooks';
import {Logout} from '../gql/__generated__/Logout';
import {Preferences, ToggleTheme} from '../gql/preferences.local';

const drawerWidth = 240;

const styles: StyleRulesCallback = (theme) => ({
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
        padding: theme.spacing.unit * 3,
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
});

const routerLink = (to: string) => {
    return (props: {}) => <Link to={to} {...props} />;
};

export const Page = withStyles(styles)(({children, classes}: React.PropsWithChildren<WithStyles<typeof styles>>) => {
    const [mobileOpen, setMobileOpen] = React.useState(false);
    const [userMenuOpen, setUserMenuOpen] = React.useState<null | HTMLElement>(null);
    const logout = useMutation<Logout>(gqlUser.Logout, {refetchQueries: [{query: gqlUser.CurrentUser}]});
    const toggleTheme = useMutation<{}>(ToggleTheme, {refetchQueries: [{query: Preferences}]});

    const drawer = (
        <div>
            <div className={classes.toolbar}>
                <HrefLink href="https://github.com/traggo" underline="none">
                    <Typography variant="h5" align="center">
                        traggo
                    </Typography>
                </HrefLink>
                <HrefLink href="https://github.com/traggo/server/releases" underline="none">
                    <Typography variant="subtitle2" align="center">
                        v1.0.0@ececece
                    </Typography>
                </HrefLink>
            </div>
            <Divider />
            <List>
                <ListItem button component={routerLink('/dashboard')}>
                    <ListItemIcon>
                        <DashboardIcon />
                    </ListItemIcon>
                    <ListItemText primary={'Dashboard'} />
                </ListItem>
            </List>
            <Divider />
            <List subheader={<ListSubheader>Timesheet</ListSubheader>}>
                <ListItem button component={routerLink('/timesheet/daily')}>
                    <ListItemIcon>
                        <ViewLintIcon />
                    </ListItemIcon>
                    <ListItemText primary="Daily" />
                </ListItem>
                <ListItem button component={routerLink('/timesheet/weekly')}>
                    <ListItemIcon>
                        <CalendarIcon />
                    </ListItemIcon>
                    <ListItemText primary="Weekly" />
                </ListItem>
            </List>
            <Divider />
            <List subheader={<ListSubheader>User</ListSubheader>}>
                <ListItem button component={routerLink('/user/settings')}>
                    <ListItemIcon>
                        <SettingsIcon />
                    </ListItemIcon>
                    <ListItemText primary="Settings" />
                </ListItem>
                <ListItem button component={routerLink('/user/devices')}>
                    <ListItemIcon>
                        <DevicesIcon />
                    </ListItemIcon>
                    <ListItemText primary="Devices" />
                </ListItem>
            </List>
            <Divider />
            <List subheader={<ListSubheader>Admin</ListSubheader>}>
                <ListItem button component={routerLink('/admin/users')}>
                    <ListItemIcon>
                        <UsersIcon />
                    </ListItemIcon>
                    <ListItemText primary="Users" />
                </ListItem>
            </List>
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
                        timesheet / daily
                    </Typography>
                    <div className={classes.grow} />
                    <div className={classes.sectionDesktop}>
                        <Button color="inherit" onClick={(e) => setUserMenuOpen(e.currentTarget)}>
                            <AccountCircle />
                            &nbsp;admin
                        </Button>
                        <IconButton color="inherit" onClick={() => toggleTheme()}>
                            <HighlightIcon />
                        </IconButton>
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
});
