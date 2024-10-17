export default function TargetList({ targets, pageId, searchString }) {
    return (
      <div>
        <h1>Page {pageId}</h1>
        <h2>Targets that begin with {searchString}</h2>
        <table border="1">
          <thead>
            <tr>
              <th>Name</th>
              <th>Info</th>
            </tr>
          </thead>
          <tbody>
            {targets.map((target) => (
              <tr key={target.id}>
                <td>{target.name}</td>
                <td>{target.info}</td>
              </tr>
            ))}
          </tbody>
        </table>
        <a href="/">Go Back Home</a>
      </div>
    );
  }