import {createRoot} from 'react-dom/client';
import {Root} from './Root';
import 'moment/locale/de';
import 'moment/locale/en-au';
import 'moment/locale/en-gb';

createRoot(document.getElementById('root')!).render(<Root />);
