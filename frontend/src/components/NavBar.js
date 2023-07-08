import Nav from 'react-bootstrap/Nav';
import Container from 'react-bootstrap/Container';
import Navbar from 'react-bootstrap/Navbar';

function NavBar() {
  return (
    <Navbar bg="light" expand="lg">
      <Container>
      <Nav className="mr-auto">
        <Nav.Link href="/watcher-list">Add plugins list</Nav.Link>
        <Nav.Link href="/plugin-manager">Plugin Manager</Nav.Link>
        <Nav.Link href="/api-stats">API stats</Nav.Link>
      </Nav>
      </Container>
    </Navbar>

  );
}

export default NavBar;