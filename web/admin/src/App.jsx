// Note: Using Babel Standalone in the shell, so we use global React/Mantine for now
// or we use ESM imports if we set type="module" in the shell script tag.
// Let's use ESM imports for cleaner "Modern" look, even if served via CDN.

import React from 'https://esm.sh/react@18.2.0';
import { createRoot } from 'https://esm.sh/react-dom@18.2.0/client';
import { 
    MantineProvider, AppShell, Navbar, Header, Text, Burger, useMantineTheme,
    NavLink, Button, Group, TextInput, PasswordInput, Paper, Title, Container, 
    Table, Badge, LoadingOverlay, Box, MediaQuery, Flex
} from 'https://esm.sh/@mantine/core@6.0.21?deps=react@18.2.0,react-dom@18.2.0,@emotion/react@11.11.1';
import * as TablerIcons from 'https://esm.sh/@tabler/icons-react@2.40.0?deps=react@18.2.0';
import axios from 'https://esm.sh/axios@1.5.0';

// --- COMPONENTS ---

const Login = ({ onLogin }) => {
    const [username, setUsername] = React.useState('');
    const [password, setPassword] = React.useState('');
    const [error, setError] = React.useState('');

    const handleSubmit = async () => {
        try {
            const res = await axios.post('/admin/v1/login', { username, password });
            if (res.data.code === 0) {
                onLogin(res.data.data.token);
            } else {
                setError(res.data.msg);
            }
        } catch (e) {
            setError('Login failed');
        }
    };

    return (
        <Container size={420} my={40}>
            <Title align="center" sx={(theme) => ({ fontFamily: `Greycliff CF, ${theme.fontFamily}`, fontWeight: 900 })}>
                Appsite Admin
            </Title>
            <Paper withBorder shadow="md" p={30} mt={30} radius="md">
                <TextInput label="Username" placeholder="admin" required value={username} onChange={(e) => setUsername(e.target.value)} />
                <PasswordInput label="Password" placeholder="password" required mt="md" value={password} onChange={(e) => setPassword(e.target.value)} />
                {error && <Text color="red" size="sm" mt="sm">{error}</Text>}
                <Button fullWidth mt="xl" onClick={handleSubmit}>Sign in</Button>
            </Paper>
        </Container>
    );
};

const UserList = ({ token }) => {
    const [data, setData] = React.useState([]);
    const [loading, setLoading] = React.useState(false);

    const fetchUsers = async () => {
        setLoading(true);
        try {
            const res = await axios.get('/admin/v1/users', { headers: { Authorization: `Bearer ${token}` } });
            if (res.data.code === 0) {
                setData(res.data.data.list);
            }
        } catch (e) { console.error(e); }
        setLoading(false);
    };

    React.useEffect(() => { fetchUsers(); }, []);

    const rows = data.map((user) => (
        <tr key={user.uid}>
            <td>{user.uid}</td>
            <td>{user.username}</td>
            <td>{user.nickname}</td>
            <td><Badge color={user.status === 'enabled' ? 'green' : 'red'}>{user.status}</Badge></td>
            <td>
                <Button variant="subtle" size="xs">Edit</Button>
            </td>
        </tr>
    ));

    return (
        <Box pos="relative">
            <LoadingOverlay visible={loading} overlayBlur={2} />
            <Group position="apart" mb="md">
                <Title order={3}>Users</Title>
                <Button onClick={fetchUsers} size="xs" variant="outline">Refresh</Button>
            </Group>
            <Paper shadow="xs" p="md">
                <Table>
                    <thead>
                        <tr>
                            <th>UID</th>
                            <th>Username</th>
                            <th>Nickname</th>
                            <th>Status</th>
                            <th />
                        </tr>
                    </thead>
                    <tbody>{rows}</tbody>
                </Table>
            </Paper>
        </Box>
    );
};

const Dashboard = () => {
    return (
        <Container>
            <Title order={2}>Dashboard</Title>
            <Text>Welcome to the Appsite Management Panel.</Text>
        </Container>
    );
};

const MainLayout = ({ token, onLogout }) => {
    const theme = useMantineTheme();
    const [opened, setOpened] = React.useState(false);
    const [page, setPage] = React.useState('dashboard');
    const [menuItems, setMenuItems] = React.useState([]);

    React.useEffect(() => {
        axios.get('/admin/v1/menu').then(res => {
            if (res.data.code === 0) setMenuItems(res.data.data);
        });
    }, []);

    const renderIcon = (iconName) => {
        const Icon = TablerIcons[iconName] || TablerIcons.IconCircle;
        return <Icon size="1rem" stroke={1.5} />;
    };

    const renderNav = (items) => items.map(item => (
        <NavLink 
            key={item.key} 
            label={item.label} 
            icon={item.icon ? renderIcon(item.icon) : null}
            active={page === item.key}
            defaultOpened={true}
            onClick={() => { if(!item.children) { setPage(item.key); setOpened(false); } }}
        >
            {item.children && renderNav(item.children)}
        </NavLink>
    ));

    const renderContent = () => {
        switch(page) {
            case 'dashboard': return <Dashboard />;
            case 'user-list': return <UserList token={token} />;
            default: return <Text>Page not found or pending implementation</Text>;
        }
    };

    return (
        <AppShell
            styles={{
                main: { background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.colors.gray[0] },
            }}
            navbarOffsetBreakpoint="sm"
            asideOffsetBreakpoint="sm"
            navbar={
                <Navbar p="md" hiddenBreakpoint="sm" hidden={!opened} width={{ sm: 200, lg: 300 }}>
                    <Navbar.Section>
                        <Text weight={700} size="sm" color="dimmed" mb="xs" transform="uppercase">Menu</Text>
                    </Navbar.Section>
                    <Navbar.Section grow mt="md">
                        {renderNav(menuItems)}
                    </Navbar.Section>
                </Navbar>
            }
            header={
                <Header height={{ base: 50, md: 70 }} p="md">
                    <div style={{ display: 'flex', alignItems: 'center', height: '100%', justifyContent: 'space-between' }}>
                        <MediaQuery largerThan="sm" styles={{ display: 'none' }}>
                            <Burger opened={opened} onClick={() => setOpened((o) => !o)} size="sm" color={theme.colors.gray[6]} mr="xl" />
                        </MediaQuery>
                        <Group>
                            <Text size="lg" weight={900} variant="gradient" gradient={{ from: 'indigo', to: 'cyan', deg: 45 }}>Appsite Admin</Text>
                        </Group>
                        <Button variant="default" onClick={onLogout}>Logout</Button>
                    </div>
                </Header>
            }
        >
            {renderContent()}
        </AppShell>
    );
};

const App = () => {
    const [token, setToken] = React.useState(localStorage.getItem('admin_token'));

    const handleLogin = (t) => {
        localStorage.setItem('admin_token', t);
        setToken(t);
    };

    const handleLogout = () => {
        localStorage.removeItem('admin_token');
        setToken(null);
    };

    return (
        <MantineProvider withGlobalStyles withNormalizeCSS theme={{ colorScheme: 'light' }}>
            {token ? <MainLayout token={token} onLogout={handleLogout} /> : <Login onLogin={handleLogin} />}
        </MantineProvider>
    );
};

const root = createRoot(document.getElementById('root'));
root.render(<App />);

export default App;
