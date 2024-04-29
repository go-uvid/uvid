import { RouterProvider, createBrowserRouter } from "react-router-dom";
import { ConfigProvider, theme } from "antd";
import "./App.css";
import { Home } from "./pages/home";

const router = createBrowserRouter([
  {
    path: "/",
    element: <Home />,
  },
]);

function App() {
  return (
    <ConfigProvider theme={{ algorithm: theme.defaultAlgorithm }}>
      <RouterProvider router={router} />
    </ConfigProvider>
  );
}

export default App;
