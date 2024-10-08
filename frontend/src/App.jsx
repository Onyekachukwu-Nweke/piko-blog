import Navbar from './components/navbar/navbar.components';
import HeroSection from './components/hero/hero.components';
import FeatureSection from './components/feature/feature.component';
import AboutUsAndMission from './components/aboutAndMission/aboutUs.component';
import CategorySection from './components/categorySection/categories.component';
import WhyWeStarted from './components/whyWeStarted/whyWeStarted.components';
import TopAuthors from './components/topAuthors/topAuthors.components';
import FeaturedIn from './components/clout/clout.components';
import Testimonials from './components/testimonials/testimonials.components';
import JoinUs from './components/joinus/joinus.component';
import Footer from './components/footer/footer.component';

function App() {

  return (
    <>
      <Navbar />
      <HeroSection />
      <FeatureSection />
      <AboutUsAndMission />
      <CategorySection />
      <WhyWeStarted />
      <TopAuthors />
      <FeaturedIn />
      <Testimonials />
      <JoinUs />
      <Footer />
    </>
  )
}

export default App
