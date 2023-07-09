
import 'bootstrap/dist/css/bootstrap.min.css';

import NavBar from './components/NavBar';

import {BrowserRouter as Router, Route, Routes} from 'react-router-dom';

import PluginManager from './pages/plugin-manager';
import PluginChanges from './pages/plugin-changes';
import AddPluginList from './pages/add-plugin-list';
import ApiStats from './pages/apiStats';

function App() {

  return (
    <>
    <Router>
      <div>
        <NavBar />
        <Routes>
          <Route path='/api-stats' element={<ApiStats/>} />
          <Route path="/plugin-manager" element={<PluginManager/>} />
          <Route path="/add-plugin-list" element={<AddPluginList/>} />
          <Route path="/plugin-changes" element={<PluginChanges/>} />
        </Routes>
      </div>
    </Router>
    </>

  );
}

export default App;