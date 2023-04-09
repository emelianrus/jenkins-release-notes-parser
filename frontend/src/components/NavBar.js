import Nav from 'react-bootstrap/Nav';
import Container from 'react-bootstrap/Container';
import Navbar from 'react-bootstrap/Navbar';

function NavBar() {
  return (
    <Navbar bg="light" expand="lg">
      <Container>
      <Nav className="mr-auto">
        <Navbar.Brand href="/" className="mr-auto">TITLE</Navbar.Brand>
        <Nav.Link href="/release-notes">Releases</Nav.Link>
        <Nav.Link href="/projects">Projects</Nav.Link>
        <Nav.Link href="/github">Github</Nav.Link>
      </Nav>
      </Container>
    </Navbar>
  );
}

export default NavBar;