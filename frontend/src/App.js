
import 'bootstrap/dist/css/bootstrap.min.css';

import NavBar from './components/NavBar';

import { BrowserRouter as Router, Route, Routes, Navigate } from 'react-router-dom';

import PluginManager from './pages/plugin-manager';
import PluginChanges from './pages/plugin-changes';
import AddPluginList from './pages/add-plugin-list';
import AddUpdatedPlugins from './pages/add-updated-plugins';
import ApiStats from './pages/apiStats';

function App() {
  const footerStyle = {
    marginTop: '20px', // Adjust the margin as needed
    backgroundColor: '#f0f0f0', // Set your desired background color
    padding: '10px', // Set padding if needed
    textAlign: 'center',
  };
  const currentYear = new Date().getFullYear();
  return (
    <>
    <Router>
      <div>
        <NavBar />
        <Routes>
          <Route path="/" element={<Navigate to ="/plugin-manager" />}/>

          <Route path='/api-stats' element={<ApiStats/>} />
          <Route path="/plugin-manager" element={<PluginManager/>} />
          <Route path="/add-plugin-list" element={<AddPluginList/>} />
          <Route path="/add-updated-plugins" element={<AddUpdatedPlugins/>} />
          <Route path="/plugin-changes" element={<PluginChanges/>} />
        </Routes>
        <footer style={footerStyle}>
          © {currentYear} Jenkins Plugin Manager
        </footer>
      </div>
    </Router>
    </>

  );
}

export default App;