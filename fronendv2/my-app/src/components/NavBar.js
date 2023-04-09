import Nav from 'react-bootstrap/Nav';
import Container from 'react-bootstrap/Container';
import Navbar from 'react-bootstrap/Navbar';

function NavBar() {
  return (
    <Navbar bg="light" expand="lg">
      <Container>
      <Nav className="mr-auto">
        <Navbar.Brand href="#home" className="mr-auto">TITLE</Navbar.Brand>
        <Nav.Link href="#home">Servers</Nav.Link>
        <Nav.Link href="#link">Projects</Nav.Link>
        <Nav.Link href="#link">Github</Nav.Link>
      </Nav>
      </Container>
    </Navbar>
  );
}

export default NavBar;