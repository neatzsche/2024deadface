export default function Home() {
    return (
      <div>
        <h1>Welcome to the Home Page</h1>
        <nav>
          <a href="/pages?page=1">Page 1</a><br />
          <a href="/pages?page=2">Page 2</a><br />
          <a href="/pages?page=3">Page 3</a><br />
          <a href="/login">Login</a>
        </nav>
      </div>
    );
  }