import logo from './logo.svg';
import './App.css';
import { useQuery, useSubscription, gql } from '@apollo/client';

const MONITOR_DATA_SUBSCRIPTION = gql`
  subscription MySubscription {
    ixp_monitoring_data_v2(order_by: {ixp: asc}) {
      city
      country
      ixp
      number_of_peers
      number_of_rib_entries
      rs_local_asn
      total_number_of_neighbors
      updated
    }
  }
`

function App() {
  // const { loading, error, data } = useQuery(MONITOR_DATA);
  const { loading, error, data } = useSubscription(MONITOR_DATA_SUBSCRIPTION);
  console.log({loading, error, data})
  if (loading) return <p>Loading...</p>;
  if (error) return <p>Error :(</p>;

  return (
    <div className="App">
      <table>
        <thead>
          <tr>
            <th>IXP</th>
            <th>City</th>
            <th>Country</th>
            <th>Updated</th>
          </tr>
        </thead>
        <tbody>
      {
        data.ixp_monitoring_data_v2.map(({ ixp, city, country, updated }) => (
        <tr key={ixp}>
          <td>{ixp}</td>
          <td>{city}</td>
          <td>{country}</td>
          <td>{updated}</td>
        </tr>
        ))
      }
        </tbody>
      </table>
    </div>
  );
}

export default App;
