package {{.Package}};

import static org.assertj.core.api.Assertions.*;

import dev.galasa.Test;
import dev.galasa.core.manager.CoreManager;
import dev.galasa.core.manager.ICoreManager;

@Test
public class SampleTest {

	// Galasa will inject a core manager into the following field
	@CoreManager
	public ICoreManager core;
	
	@Test
	public void simpleSampleTest() {
		assertThat(core).isNotNull();
	}

}
