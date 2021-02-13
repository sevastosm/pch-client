import logo from './logo.svg';
import './App.css';
import { useQuery, useSubscription, gql } from '@apollo/client';

const MONITOR_DATA_SUBSCRIPTION = gql`
  subscription MySubscription {
    ixp_server_data(order_by: {ixp: asc}) {
      ixp
      city
      country
      protocol
      number_of_peers
      number_of_rib_entries
      rs_local_asn
      total_number_of_neighbors
      updated_at
    }
  }
`

function App() {
  const { loading, error, data } = useSubscription(MONITOR_DATA_SUBSCRIPTION);
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
            <th>Protocol</th>
            <th>Number of peers</th>
            <th>Number of RIB entries</th>
            <th>Local ASN</th>
            <th>Number of Neighbors</th>
            <th>Updated</th>
          </tr>
        </thead>
        <tbody>
      {
        data.ixp_server_data.map(({ ixp, city, country, protocol, number_of_peers, number_of_rib_entries, rs_local_asn, total_number_of_neighbors, updated_at }) => (
        <tr key={ixp}>
          <td>{ixp}</td>
          <td>{city}</td>
          <td>{country}</td>
          <td>{protocol}</td>
          <td>{number_of_peers}</td>
          <td>{number_of_rib_entries}</td>
          <td>{rs_local_asn}</td>
          <td>{total_number_of_neighbors}</td>
          <td>{updated_at}</td>
        </tr>
        ))
      }
        </tbody>
      </table>
    </div>
  );
}

export default App;
